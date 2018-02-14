-- Copyright 2018 The Chromium Authors. All Rights Reserved.
--
-- Use of this source code is governed by a BSD-style
-- license that can be found in the LICENSE file or at
-- https://developers.google.com/open-source/licenses/bsd


DROP PROCEDURE IF EXISTS BackfillIssueSnapshotsCcs;
DROP PROCEDURE IF EXISTS BackfillIssueSnapshotsComponents;
DROP PROCEDURE IF EXISTS BackfillIssueSnapshotsLabels;
DROP PROCEDURE IF EXISTS BackfillIssueSnapshotsChunk;
DROP PROCEDURE IF EXISTS BackfillIssueSnapshotsManyChunks;

delimiter //

CREATE PROCEDURE BackfillIssueSnapshotsLabels(IN c_issue_id INT, IN c_issuesnapshot_id INT)
BEGIN

  DECLARE done INT DEFAULT FALSE;

  DECLARE c_label_id INT;

  DECLARE curs CURSOR FOR
    SELECT label_id
    FROM Issue2Label
    WHERE issue_id = c_issue_id;

  DECLARE CONTINUE HANDLER FOR NOT FOUND SET done = TRUE;
  OPEN curs;

  label_loop: LOOP
    FETCH curs INTO c_label_id;
    IF done THEN
      LEAVE label_loop;
    END IF;

    INSERT INTO IssueSnapshot2Label
      (issuesnapshot_id, label_id)
      VALUES
      (c_issuesnapshot_id, c_label_id);

  END LOOP;

  CLOSE curs;
END;


CREATE PROCEDURE BackfillIssueSnapshotsCcs(IN c_issue_id INT, IN c_issuesnapshot_id INT)
BEGIN

  DECLARE done INT DEFAULT FALSE;

  DECLARE c_cc_id INT UNSIGNED;

  DECLARE curs CURSOR FOR
    SELECT cc_id
    FROM Issue2Cc
    WHERE issue_id = c_issue_id;

  DECLARE CONTINUE HANDLER FOR NOT FOUND SET done = TRUE;
  OPEN curs;

  cc_loop: LOOP
    FETCH curs INTO c_cc_id;
    IF done THEN
      LEAVE cc_loop;
    END IF;

    INSERT INTO IssueSnapshot2Cc
      (issuesnapshot_id, cc_id)
      VALUES
      (c_issuesnapshot_id, c_cc_id);

  END LOOP;

  CLOSE curs;
END;


CREATE PROCEDURE BackfillIssueSnapshotsComponents(IN c_issue_id INT, IN c_issuesnapshot_id INT)
BEGIN

  DECLARE done INT DEFAULT FALSE;

  DECLARE c_component_id INT;

  DECLARE curs CURSOR FOR
    SELECT component_id
    FROM Issue2Component
    WHERE issue_id = c_issue_id;

  DECLARE CONTINUE HANDLER FOR NOT FOUND SET done = TRUE;
  OPEN curs;

  component_loop: LOOP
    FETCH curs INTO c_component_id;
    IF done THEN
      LEAVE component_loop;
    END IF;

    INSERT INTO IssueSnapshot2Component
      (issuesnapshot_id, component_id)
      VALUES
      (c_issuesnapshot_id, c_component_id);

  END LOOP;

  CLOSE curs;
END;


CREATE PROCEDURE BackfillIssueSnapshotsChunk(IN chunk_size INT UNSIGNED,
  IN chunk_offset INT UNSIGNED)
BEGIN

  DECLARE done TINYINT DEFAULT FALSE;
  DECLARE is_duplicate TINYINT DEFAULT FALSE;

  DECLARE c_issue_id INT;
  DECLARE c_issue_shard INT;
  DECLARE c_issue_project_id INT;
  DECLARE c_issue_local_id INT;
  DECLARE c_issue_status_id INT;
  DECLARE c_issue_opened INT;
  DECLARE c_issue_closed INT;
  DECLARE c_issue_is_open BOOLEAN;
  DECLARE c_reporter_id INT UNSIGNED;
  DECLARE c_owner_id INT UNSIGNED;
  DECLARE c_issuesnapshot_id INT;

  DECLARE curs CURSOR FOR
    SELECT i.id, i.shard, i.project_id, i.local_id, i.status_id, i.opened,
      CASE sd.means_open WHEN true THEN (
        -- If a snapshot for this Issue already exists, make the new snapshot's
        -- period_end the period_start of the existing snapshot.
        SELECT IFNULL((
            SELECT period_start FROM IssueSnapshot WHERE issue_id = i.id LIMIT 1
        ), 4294967295)
      ) ELSE i.closed END,
      sd.means_open,
      i.reporter_id, i.owner_id
    FROM Issue i
    JOIN StatusDef sd ON i.status_id = sd.id
    WHERE i.id >= chunk_offset AND i.id < chunk_offset + chunk_size;

  DECLARE CONTINUE HANDLER FOR NOT FOUND SET done = TRUE;

  -- Handle duplicates. Using this along with UNIQUE indexes is how we make
  -- the backfill idempotent.
  DECLARE CONTINUE HANDLER FOR SQLSTATE '23000' BEGIN
    SET is_duplicate = TRUE;
  END;

  OPEN curs;

  issue_loop: LOOP
    SET is_duplicate = FALSE;

    FETCH curs INTO c_issue_id, c_issue_shard, c_issue_project_id,
      c_issue_local_id, c_issue_status_id, c_issue_opened, c_issue_closed,
      c_issue_is_open, c_reporter_id, c_owner_id;
    IF is_duplicate THEN
      ITERATE issue_loop;
    END IF;
    IF done THEN
      LEAVE issue_loop;
    END IF;

    -- Indicate progress.
    SELECT c_issue_id AS 'Processing Issue:';

    INSERT INTO IssueSnapshot
    (issue_id, shard, project_id, local_id, status_id, period_start,
    period_end, is_open, reporter_id, owner_id)
    VALUES
    (c_issue_id, c_issue_shard, c_issue_project_id,
    c_issue_local_id, c_issue_status_id, c_issue_opened,
    c_issue_closed, c_issue_is_open, c_reporter_id, c_owner_id);

    SET c_issuesnapshot_id = LAST_INSERT_ID();
    SELECT c_issuesnapshot_id AS 'IssueSnapshot ID:';

    -- Backfill labels.
    CALL BackfillIssueSnapshotsLabels(c_issue_id, c_issuesnapshot_id);
    CALL BackfillIssueSnapshotsCcs(c_issue_id, c_issuesnapshot_id);
    CALL BackfillIssueSnapshotsComponents(c_issue_id, c_issuesnapshot_id);

  END LOOP;

  CLOSE curs;
END;


CREATE PROCEDURE BackfillIssueSnapshotsManyChunks(
  IN num_chunks SMALLINT UNSIGNED,
  IN chunk_size SMALLINT UNSIGNED)
BEGIN
  DECLARE chunk_i INT DEFAULT 0;

  -- Handle no results found ("cursor is not open")
  DECLARE CONTINUE HANDLER FOR SQLSTATE '24000' BEGIN END;

  WHILE chunk_i < num_chunks DO

    SELECT CONCAT('Backfilling chunk ', chunk_i + 1, ' of ', num_chunks) AS '';
    SELECT chunk_size, chunk_i, chunk_i * chunk_size AS 'chunk offset';

    CALL BackfillIssueSnapshotsChunk(chunk_size, chunk_i * chunk_size);

    SELECT SLEEP(1);

    SET chunk_i = chunk_i + 1;
  END WHILE;
END;


//


delimiter ;
