package bbutil

import (
	"time"

	"golang.org/x/net/context"
	"google.golang.org/api/googleapi"

	"github.com/luci/luci-go/common/api/buildbucket/buildbucket/v1"
	"github.com/luci/luci-go/common/errors"
	"github.com/luci/luci-go/common/logging"
	"github.com/luci/luci-go/common/retry"
)

// Search searches for builds and sends them to the builds channel
// until the context is cancelled or there are no more builds.
// The order from most-recently-created to least-recently-created.
//
// If minCreationDate is not zero, stops searching when reaches a build created
// before minCreationDate.
//
// Returns nil only if the search results are exhausted. Otherwise it returns
// c.Err() or a fatal error.
//
// On transient errors, logs them and retries requests with an exponential
// back-off.
func Search(c context.Context, req *buildbucket.SearchCall, minCreationDate time.Time, builds chan<- *buildbucket.ApiCommonBuildMessage) error {
	for {
		var batch *buildbucket.ApiSearchResponseMessage
		err := retry.Retry(c, retry.TransientOnly(retry.Default),
			func() error {
				reqCtx, _ := context.WithTimeout(c, time.Minute)
				res, err := req.Context(reqCtx).Do()
				switch apiErr, _ := err.(*googleapi.Error); {
				case apiErr != nil && apiErr.Code >= 500:
					return errors.WrapTransient(err)
				case err == context.DeadlineExceeded && c.Err() == nil:
					return errors.WrapTransient(err)
				case err != nil:
					return err
				case res.Error != nil:
					return errors.New(res.Error.Message)
				default:
					batch = res
					return nil
				}
			},
			func(err error, wait time.Duration) {
				logging.WithError(err).Warningf(c, "transient RPC error while searching builds; will retry in %s", wait)
			})
		if err != nil {
			return err
		}

		for _, b := range batch.Builds {
			if !minCreationDate.IsZero() && ParseTimestamp(b.CreatedTs).Before(minCreationDate) {
				// search results are always ordered newest to oldest.
				return nil
			}
			select {
			case <-c.Done():
				return c.Err()
			case builds <- b:
			}
		}

		if len(batch.Builds) == 0 || batch.NextCursor == "" {
			break
		}
		req.StartCursor(batch.NextCursor)
	}

	return nil
}

// SearchAll is similar to Search, but returns builds as a slice.
func SearchAll(c context.Context, req *buildbucket.SearchCall, minCreationDate time.Time) ([]*buildbucket.ApiCommonBuildMessage, error) {
	ch := make(chan *buildbucket.ApiCommonBuildMessage)
	var err error
	go func() {
		defer close(ch)
		err = Search(c, req, minCreationDate, ch)
	}()

	var builds []*buildbucket.ApiCommonBuildMessage
	for b := range ch {
		builds = append(builds, b)
	}
	return builds, err
}
