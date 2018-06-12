# Generated by the protocol buffer compiler.  DO NOT EDIT!
# source: api/api_proto/issues.proto

import sys
_b=sys.version_info[0]<3 and (lambda x:x) or (lambda x:x.encode('latin1'))
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from google.protobuf import reflection as _reflection
from google.protobuf import symbol_database as _symbol_database
from google.protobuf import descriptor_pb2
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()


from google.protobuf import empty_pb2 as google_dot_protobuf_dot_empty__pb2
from api.api_proto import common_pb2 as api_dot_api__proto_dot_common__pb2
from api.api_proto import issue_objects_pb2 as api_dot_api__proto_dot_issue__objects__pb2


DESCRIPTOR = _descriptor.FileDescriptor(
  name='api/api_proto/issues.proto',
  package='monorail',
  syntax='proto3',
  serialized_pb=_b('\n\x1a\x61pi/api_proto/issues.proto\x12\x08monorail\x1a\x1bgoogle/protobuf/empty.proto\x1a\x1a\x61pi/api_proto/common.proto\x1a!api/api_proto/issue_objects.proto\"q\n\x12\x43reateIssueRequest\x12%\n\x05trace\x18\x01 \x01(\x0b\x32\x16.monorail.RequestTrace\x12\x14\n\x0cproject_name\x18\x02 \x01(\t\x12\x1e\n\x05issue\x18\x03 \x01(\x0b\x32\x0f.monorail.Issue\"_\n\x0fGetIssueRequest\x12%\n\x05trace\x18\x01 \x01(\x0b\x32\x16.monorail.RequestTrace\x12%\n\tissue_ref\x18\x02 \x01(\x0b\x32\x12.monorail.IssueRef\"/\n\rIssueResponse\x12\x1e\n\x05issue\x18\x01 \x01(\x0b\x32\x0f.monorail.Issue\"c\n\x13ListCommentsRequest\x12%\n\x05trace\x18\x01 \x01(\x0b\x32\x16.monorail.RequestTrace\x12%\n\tissue_ref\x18\x02 \x01(\x0b\x32\x12.monorail.IssueRef\";\n\x14ListCommentsResponse\x12#\n\x08\x63omments\x18\x01 \x03(\x0b\x32\x11.monorail.Comment\"\x8e\x01\n\x19\x44\x65leteIssueCommentRequest\x12%\n\x05trace\x18\x01 \x01(\x0b\x32\x16.monorail.RequestTrace\x12\x14\n\x0cproject_name\x18\x02 \x01(\t\x12\x10\n\x08local_id\x18\x03 \x01(\x03\x12\x12\n\ncomment_id\x18\x04 \x01(\x03\x12\x0e\n\x06\x64\x65lete\x18\x05 \x01(\x08\"\xea\x01\n\x15UpdateApprovalRequest\x12%\n\x05trace\x18\x01 \x01(\x0b\x32\x16.monorail.RequestTrace\x12%\n\tissue_ref\x18\x02 \x01(\x0b\x32\x12.monorail.IssueRef\x12%\n\tfield_ref\x18\x03 \x01(\x0b\x32\x12.monorail.FieldRef\x12/\n\x0e\x61pproval_delta\x18\x04 \x01(\x0b\x32\x17.monorail.ApprovalDelta\x12\x17\n\x0f\x63omment_content\x18\x05 \x01(\t\x12\x12\n\nsend_email\x18\x06 \x01(\x08\">\n\x16UpdateApprovalResponse\x12$\n\x08\x61pproval\x18\x01 \x01(\x0b\x32\x12.monorail.Approval2\x8f\x03\n\x06Issues\x12\x46\n\x0b\x43reateIssue\x12\x1c.monorail.CreateIssueRequest\x1a\x17.monorail.IssueResponse\"\x00\x12@\n\x08GetIssue\x12\x19.monorail.GetIssueRequest\x1a\x17.monorail.IssueResponse\"\x00\x12O\n\x0cListComments\x12\x1d.monorail.ListCommentsRequest\x1a\x1e.monorail.ListCommentsResponse\"\x00\x12S\n\x12\x44\x65leteIssueComment\x12#.monorail.DeleteIssueCommentRequest\x1a\x16.google.protobuf.Empty\"\x00\x12U\n\x0eUpdateApproval\x12\x1f.monorail.UpdateApprovalRequest\x1a .monorail.UpdateApprovalResponse\"\x00\x62\x06proto3')
  ,
  dependencies=[google_dot_protobuf_dot_empty__pb2.DESCRIPTOR,api_dot_api__proto_dot_common__pb2.DESCRIPTOR,api_dot_api__proto_dot_issue__objects__pb2.DESCRIPTOR,])
_sym_db.RegisterFileDescriptor(DESCRIPTOR)




_CREATEISSUEREQUEST = _descriptor.Descriptor(
  name='CreateIssueRequest',
  full_name='monorail.CreateIssueRequest',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='trace', full_name='monorail.CreateIssueRequest.trace', index=0,
      number=1, type=11, cpp_type=10, label=1,
      has_default_value=False, default_value=None,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      options=None),
    _descriptor.FieldDescriptor(
      name='project_name', full_name='monorail.CreateIssueRequest.project_name', index=1,
      number=2, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=_b("").decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      options=None),
    _descriptor.FieldDescriptor(
      name='issue', full_name='monorail.CreateIssueRequest.issue', index=2,
      number=3, type=11, cpp_type=10, label=1,
      has_default_value=False, default_value=None,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      options=None),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=132,
  serialized_end=245,
)


_GETISSUEREQUEST = _descriptor.Descriptor(
  name='GetIssueRequest',
  full_name='monorail.GetIssueRequest',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='trace', full_name='monorail.GetIssueRequest.trace', index=0,
      number=1, type=11, cpp_type=10, label=1,
      has_default_value=False, default_value=None,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      options=None),
    _descriptor.FieldDescriptor(
      name='issue_ref', full_name='monorail.GetIssueRequest.issue_ref', index=1,
      number=2, type=11, cpp_type=10, label=1,
      has_default_value=False, default_value=None,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      options=None),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=247,
  serialized_end=342,
)


_ISSUERESPONSE = _descriptor.Descriptor(
  name='IssueResponse',
  full_name='monorail.IssueResponse',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='issue', full_name='monorail.IssueResponse.issue', index=0,
      number=1, type=11, cpp_type=10, label=1,
      has_default_value=False, default_value=None,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      options=None),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=344,
  serialized_end=391,
)


_LISTCOMMENTSREQUEST = _descriptor.Descriptor(
  name='ListCommentsRequest',
  full_name='monorail.ListCommentsRequest',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='trace', full_name='monorail.ListCommentsRequest.trace', index=0,
      number=1, type=11, cpp_type=10, label=1,
      has_default_value=False, default_value=None,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      options=None),
    _descriptor.FieldDescriptor(
      name='issue_ref', full_name='monorail.ListCommentsRequest.issue_ref', index=1,
      number=2, type=11, cpp_type=10, label=1,
      has_default_value=False, default_value=None,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      options=None),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=393,
  serialized_end=492,
)


_LISTCOMMENTSRESPONSE = _descriptor.Descriptor(
  name='ListCommentsResponse',
  full_name='monorail.ListCommentsResponse',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='comments', full_name='monorail.ListCommentsResponse.comments', index=0,
      number=1, type=11, cpp_type=10, label=3,
      has_default_value=False, default_value=[],
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      options=None),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=494,
  serialized_end=553,
)


_DELETEISSUECOMMENTREQUEST = _descriptor.Descriptor(
  name='DeleteIssueCommentRequest',
  full_name='monorail.DeleteIssueCommentRequest',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='trace', full_name='monorail.DeleteIssueCommentRequest.trace', index=0,
      number=1, type=11, cpp_type=10, label=1,
      has_default_value=False, default_value=None,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      options=None),
    _descriptor.FieldDescriptor(
      name='project_name', full_name='monorail.DeleteIssueCommentRequest.project_name', index=1,
      number=2, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=_b("").decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      options=None),
    _descriptor.FieldDescriptor(
      name='local_id', full_name='monorail.DeleteIssueCommentRequest.local_id', index=2,
      number=3, type=3, cpp_type=2, label=1,
      has_default_value=False, default_value=0,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      options=None),
    _descriptor.FieldDescriptor(
      name='comment_id', full_name='monorail.DeleteIssueCommentRequest.comment_id', index=3,
      number=4, type=3, cpp_type=2, label=1,
      has_default_value=False, default_value=0,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      options=None),
    _descriptor.FieldDescriptor(
      name='delete', full_name='monorail.DeleteIssueCommentRequest.delete', index=4,
      number=5, type=8, cpp_type=7, label=1,
      has_default_value=False, default_value=False,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      options=None),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=556,
  serialized_end=698,
)


_UPDATEAPPROVALREQUEST = _descriptor.Descriptor(
  name='UpdateApprovalRequest',
  full_name='monorail.UpdateApprovalRequest',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='trace', full_name='monorail.UpdateApprovalRequest.trace', index=0,
      number=1, type=11, cpp_type=10, label=1,
      has_default_value=False, default_value=None,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      options=None),
    _descriptor.FieldDescriptor(
      name='issue_ref', full_name='monorail.UpdateApprovalRequest.issue_ref', index=1,
      number=2, type=11, cpp_type=10, label=1,
      has_default_value=False, default_value=None,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      options=None),
    _descriptor.FieldDescriptor(
      name='field_ref', full_name='monorail.UpdateApprovalRequest.field_ref', index=2,
      number=3, type=11, cpp_type=10, label=1,
      has_default_value=False, default_value=None,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      options=None),
    _descriptor.FieldDescriptor(
      name='approval_delta', full_name='monorail.UpdateApprovalRequest.approval_delta', index=3,
      number=4, type=11, cpp_type=10, label=1,
      has_default_value=False, default_value=None,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      options=None),
    _descriptor.FieldDescriptor(
      name='comment_content', full_name='monorail.UpdateApprovalRequest.comment_content', index=4,
      number=5, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=_b("").decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      options=None),
    _descriptor.FieldDescriptor(
      name='send_email', full_name='monorail.UpdateApprovalRequest.send_email', index=5,
      number=6, type=8, cpp_type=7, label=1,
      has_default_value=False, default_value=False,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      options=None),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=701,
  serialized_end=935,
)


_UPDATEAPPROVALRESPONSE = _descriptor.Descriptor(
  name='UpdateApprovalResponse',
  full_name='monorail.UpdateApprovalResponse',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='approval', full_name='monorail.UpdateApprovalResponse.approval', index=0,
      number=1, type=11, cpp_type=10, label=1,
      has_default_value=False, default_value=None,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      options=None),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=937,
  serialized_end=999,
)

_CREATEISSUEREQUEST.fields_by_name['trace'].message_type = api_dot_api__proto_dot_common__pb2._REQUESTTRACE
_CREATEISSUEREQUEST.fields_by_name['issue'].message_type = api_dot_api__proto_dot_issue__objects__pb2._ISSUE
_GETISSUEREQUEST.fields_by_name['trace'].message_type = api_dot_api__proto_dot_common__pb2._REQUESTTRACE
_GETISSUEREQUEST.fields_by_name['issue_ref'].message_type = api_dot_api__proto_dot_common__pb2._ISSUEREF
_ISSUERESPONSE.fields_by_name['issue'].message_type = api_dot_api__proto_dot_issue__objects__pb2._ISSUE
_LISTCOMMENTSREQUEST.fields_by_name['trace'].message_type = api_dot_api__proto_dot_common__pb2._REQUESTTRACE
_LISTCOMMENTSREQUEST.fields_by_name['issue_ref'].message_type = api_dot_api__proto_dot_common__pb2._ISSUEREF
_LISTCOMMENTSRESPONSE.fields_by_name['comments'].message_type = api_dot_api__proto_dot_issue__objects__pb2._COMMENT
_DELETEISSUECOMMENTREQUEST.fields_by_name['trace'].message_type = api_dot_api__proto_dot_common__pb2._REQUESTTRACE
_UPDATEAPPROVALREQUEST.fields_by_name['trace'].message_type = api_dot_api__proto_dot_common__pb2._REQUESTTRACE
_UPDATEAPPROVALREQUEST.fields_by_name['issue_ref'].message_type = api_dot_api__proto_dot_common__pb2._ISSUEREF
_UPDATEAPPROVALREQUEST.fields_by_name['field_ref'].message_type = api_dot_api__proto_dot_common__pb2._FIELDREF
_UPDATEAPPROVALREQUEST.fields_by_name['approval_delta'].message_type = api_dot_api__proto_dot_issue__objects__pb2._APPROVALDELTA
_UPDATEAPPROVALRESPONSE.fields_by_name['approval'].message_type = api_dot_api__proto_dot_issue__objects__pb2._APPROVAL
DESCRIPTOR.message_types_by_name['CreateIssueRequest'] = _CREATEISSUEREQUEST
DESCRIPTOR.message_types_by_name['GetIssueRequest'] = _GETISSUEREQUEST
DESCRIPTOR.message_types_by_name['IssueResponse'] = _ISSUERESPONSE
DESCRIPTOR.message_types_by_name['ListCommentsRequest'] = _LISTCOMMENTSREQUEST
DESCRIPTOR.message_types_by_name['ListCommentsResponse'] = _LISTCOMMENTSRESPONSE
DESCRIPTOR.message_types_by_name['DeleteIssueCommentRequest'] = _DELETEISSUECOMMENTREQUEST
DESCRIPTOR.message_types_by_name['UpdateApprovalRequest'] = _UPDATEAPPROVALREQUEST
DESCRIPTOR.message_types_by_name['UpdateApprovalResponse'] = _UPDATEAPPROVALRESPONSE

CreateIssueRequest = _reflection.GeneratedProtocolMessageType('CreateIssueRequest', (_message.Message,), dict(
  DESCRIPTOR = _CREATEISSUEREQUEST,
  __module__ = 'api.api_proto.issues_pb2'
  # @@protoc_insertion_point(class_scope:monorail.CreateIssueRequest)
  ))
_sym_db.RegisterMessage(CreateIssueRequest)

GetIssueRequest = _reflection.GeneratedProtocolMessageType('GetIssueRequest', (_message.Message,), dict(
  DESCRIPTOR = _GETISSUEREQUEST,
  __module__ = 'api.api_proto.issues_pb2'
  # @@protoc_insertion_point(class_scope:monorail.GetIssueRequest)
  ))
_sym_db.RegisterMessage(GetIssueRequest)

IssueResponse = _reflection.GeneratedProtocolMessageType('IssueResponse', (_message.Message,), dict(
  DESCRIPTOR = _ISSUERESPONSE,
  __module__ = 'api.api_proto.issues_pb2'
  # @@protoc_insertion_point(class_scope:monorail.IssueResponse)
  ))
_sym_db.RegisterMessage(IssueResponse)

ListCommentsRequest = _reflection.GeneratedProtocolMessageType('ListCommentsRequest', (_message.Message,), dict(
  DESCRIPTOR = _LISTCOMMENTSREQUEST,
  __module__ = 'api.api_proto.issues_pb2'
  # @@protoc_insertion_point(class_scope:monorail.ListCommentsRequest)
  ))
_sym_db.RegisterMessage(ListCommentsRequest)

ListCommentsResponse = _reflection.GeneratedProtocolMessageType('ListCommentsResponse', (_message.Message,), dict(
  DESCRIPTOR = _LISTCOMMENTSRESPONSE,
  __module__ = 'api.api_proto.issues_pb2'
  # @@protoc_insertion_point(class_scope:monorail.ListCommentsResponse)
  ))
_sym_db.RegisterMessage(ListCommentsResponse)

DeleteIssueCommentRequest = _reflection.GeneratedProtocolMessageType('DeleteIssueCommentRequest', (_message.Message,), dict(
  DESCRIPTOR = _DELETEISSUECOMMENTREQUEST,
  __module__ = 'api.api_proto.issues_pb2'
  # @@protoc_insertion_point(class_scope:monorail.DeleteIssueCommentRequest)
  ))
_sym_db.RegisterMessage(DeleteIssueCommentRequest)

UpdateApprovalRequest = _reflection.GeneratedProtocolMessageType('UpdateApprovalRequest', (_message.Message,), dict(
  DESCRIPTOR = _UPDATEAPPROVALREQUEST,
  __module__ = 'api.api_proto.issues_pb2'
  # @@protoc_insertion_point(class_scope:monorail.UpdateApprovalRequest)
  ))
_sym_db.RegisterMessage(UpdateApprovalRequest)

UpdateApprovalResponse = _reflection.GeneratedProtocolMessageType('UpdateApprovalResponse', (_message.Message,), dict(
  DESCRIPTOR = _UPDATEAPPROVALRESPONSE,
  __module__ = 'api.api_proto.issues_pb2'
  # @@protoc_insertion_point(class_scope:monorail.UpdateApprovalResponse)
  ))
_sym_db.RegisterMessage(UpdateApprovalResponse)


# @@protoc_insertion_point(module_scope)
