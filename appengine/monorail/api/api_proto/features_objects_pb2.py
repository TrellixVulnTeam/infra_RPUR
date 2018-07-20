# Generated by the protocol buffer compiler.  DO NOT EDIT!
# source: api/api_proto/features_objects.proto

import sys
_b=sys.version_info[0]<3 and (lambda x:x) or (lambda x:x.encode('latin1'))
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from google.protobuf import reflection as _reflection
from google.protobuf import symbol_database as _symbol_database
from google.protobuf import descriptor_pb2
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()


from api.api_proto import common_pb2 as api_dot_api__proto_dot_common__pb2


DESCRIPTOR = _descriptor.FileDescriptor(
  name='api/api_proto/features_objects.proto',
  package='monorail',
  syntax='proto3',
  serialized_pb=_b('\n$api/api_proto/features_objects.proto\x12\x08monorail\x1a\x1a\x61pi/api_proto/common.proto\"c\n\x07Hotlist\x12$\n\towner_ref\x18\x01 \x01(\x0b\x32\x11.monorail.UserRef\x12\x0c\n\x04name\x18\x02 \x01(\t\x12\x0f\n\x07summary\x18\x03 \x01(\t\x12\x13\n\x0b\x64\x65scription\x18\x04 \x01(\tb\x06proto3')
  ,
  dependencies=[api_dot_api__proto_dot_common__pb2.DESCRIPTOR,])
_sym_db.RegisterFileDescriptor(DESCRIPTOR)




_HOTLIST = _descriptor.Descriptor(
  name='Hotlist',
  full_name='monorail.Hotlist',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='owner_ref', full_name='monorail.Hotlist.owner_ref', index=0,
      number=1, type=11, cpp_type=10, label=1,
      has_default_value=False, default_value=None,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      options=None),
    _descriptor.FieldDescriptor(
      name='name', full_name='monorail.Hotlist.name', index=1,
      number=2, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=_b("").decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      options=None),
    _descriptor.FieldDescriptor(
      name='summary', full_name='monorail.Hotlist.summary', index=2,
      number=3, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=_b("").decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      options=None),
    _descriptor.FieldDescriptor(
      name='description', full_name='monorail.Hotlist.description', index=3,
      number=4, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=_b("").decode('utf-8'),
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
  serialized_start=78,
  serialized_end=177,
)

_HOTLIST.fields_by_name['owner_ref'].message_type = api_dot_api__proto_dot_common__pb2._USERREF
DESCRIPTOR.message_types_by_name['Hotlist'] = _HOTLIST

Hotlist = _reflection.GeneratedProtocolMessageType('Hotlist', (_message.Message,), dict(
  DESCRIPTOR = _HOTLIST,
  __module__ = 'api.api_proto.features_objects_pb2'
  # @@protoc_insertion_point(class_scope:monorail.Hotlist)
  ))
_sym_db.RegisterMessage(Hotlist)


# @@protoc_insertion_point(module_scope)
