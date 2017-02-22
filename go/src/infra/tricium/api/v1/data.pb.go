// Code generated by protoc-gen-go.
// source: infra/tricium/api/v1/data.proto
// DO NOT EDIT!

package tricium

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// Available data types should be listed in this enum and have a
// corresponding nested message with a mandatory platforms fields,
// see GitFileDetails for field details.
type Data_Type int32

const (
	Data_NONE             Data_Type = 0
	Data_GIT_FILE_DETAILS Data_Type = 1
	Data_FILES            Data_Type = 2
	Data_CLANG_DETAILS    Data_Type = 3
	Data_RESULTS          Data_Type = 4
)

var Data_Type_name = map[int32]string{
	0: "NONE",
	1: "GIT_FILE_DETAILS",
	2: "FILES",
	3: "CLANG_DETAILS",
	4: "RESULTS",
}
var Data_Type_value = map[string]int32{
	"NONE":             0,
	"GIT_FILE_DETAILS": 1,
	"FILES":            2,
	"CLANG_DETAILS":    3,
	"RESULTS":          4,
}

func (x Data_Type) String() string {
	return proto.EnumName(Data_Type_name, int32(x))
}
func (Data_Type) EnumDescriptor() ([]byte, []int) { return fileDescriptor1, []int{0, 0} }

// Tricium data types.
//
// Any data type provided or needed by a Tricium analyzer.
type Data struct {
}

func (m *Data) Reset()                    { *m = Data{} }
func (m *Data) String() string            { return proto.CompactTextString(m) }
func (*Data) ProtoMessage()               {}
func (*Data) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{0} }

// Details for supported types, specifically whether a type is tied to
// a platform.
//
// These type details are used to resolve data dependencies when
// generating workflows.
type Data_TypeDetails struct {
	Type               Data_Type `protobuf:"varint,1,opt,name=type,enum=tricium.Data_Type" json:"type,omitempty"`
	IsPlatformSpecific bool      `protobuf:"varint,2,opt,name=is_platform_specific,json=isPlatformSpecific" json:"is_platform_specific,omitempty"`
}

func (m *Data_TypeDetails) Reset()                    { *m = Data_TypeDetails{} }
func (m *Data_TypeDetails) String() string            { return proto.CompactTextString(m) }
func (*Data_TypeDetails) ProtoMessage()               {}
func (*Data_TypeDetails) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{0, 0} }

func (m *Data_TypeDetails) GetType() Data_Type {
	if m != nil {
		return m.Type
	}
	return Data_NONE
}

func (m *Data_TypeDetails) GetIsPlatformSpecific() bool {
	if m != nil {
		return m.IsPlatformSpecific
	}
	return false
}

// Details for retrieval of file content from a Git repository.
//
// ISOLATED PATH: tricium/data/git_file_details.json
type Data_GitFileDetails struct {
	// The platforms this data is tied to encoded as a bitmap.
	//
	// The bit number for each platform should correspond to the enum
	// position number of the same platform in the Platform.Name enum.
	//
	// This includes the ANY platform, encoded as zero, which should
	// be used for any data that is not platform-specific.
	Platforms  int64    `protobuf:"varint,1,opt,name=platforms" json:"platforms,omitempty"`
	Repository string   `protobuf:"bytes,2,opt,name=repository" json:"repository,omitempty"`
	Ref        string   `protobuf:"bytes,3,opt,name=ref" json:"ref,omitempty"`
	Path       []string `protobuf:"bytes,4,rep,name=path" json:"path,omitempty"`
}

func (m *Data_GitFileDetails) Reset()                    { *m = Data_GitFileDetails{} }
func (m *Data_GitFileDetails) String() string            { return proto.CompactTextString(m) }
func (*Data_GitFileDetails) ProtoMessage()               {}
func (*Data_GitFileDetails) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{0, 1} }

func (m *Data_GitFileDetails) GetPlatforms() int64 {
	if m != nil {
		return m.Platforms
	}
	return 0
}

func (m *Data_GitFileDetails) GetRepository() string {
	if m != nil {
		return m.Repository
	}
	return ""
}

func (m *Data_GitFileDetails) GetRef() string {
	if m != nil {
		return m.Ref
	}
	return ""
}

func (m *Data_GitFileDetails) GetPath() []string {
	if m != nil {
		return m.Path
	}
	return nil
}

// List of paths included in the isolated input.
//
// Files in the isolate should be laid out with the same file system
// structure as in the repository, with the root of the isolate input mapped
// to the root of the repository.
//
// ISOLATED PATH: tricium/data/files.json
type Data_Files struct {
	Platforms int64 `protobuf:"varint,1,opt,name=platforms" json:"platforms,omitempty"`
	// Path to files from the root of the isolated input.
	Path []string `protobuf:"bytes,2,rep,name=path" json:"path,omitempty"`
}

func (m *Data_Files) Reset()                    { *m = Data_Files{} }
func (m *Data_Files) String() string            { return proto.CompactTextString(m) }
func (*Data_Files) ProtoMessage()               {}
func (*Data_Files) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{0, 2} }

func (m *Data_Files) GetPlatforms() int64 {
	if m != nil {
		return m.Platforms
	}
	return 0
}

func (m *Data_Files) GetPath() []string {
	if m != nil {
		return m.Path
	}
	return nil
}

// Details needed to replay a clang compilation.
//
// Any included path should correspond to a file in the isolated input.
//
// ISOLATED PATH: tricium/data/clang_details.json
type Data_ClangDetails struct {
	Platforms int64 `protobuf:"varint,1,opt,name=platforms" json:"platforms,omitempty"`
	// Path to the compilation database. Typically, in the build root.
	CompilationDb string `protobuf:"bytes,2,opt,name=compilation_db,json=compilationDb" json:"compilation_db,omitempty"`
	// Paths to files needed to compile cpp files to analyze.
	CompDepPath []string `protobuf:"bytes,3,rep,name=comp_dep_path,json=compDepPath" json:"comp_dep_path,omitempty"`
}

func (m *Data_ClangDetails) Reset()                    { *m = Data_ClangDetails{} }
func (m *Data_ClangDetails) String() string            { return proto.CompactTextString(m) }
func (*Data_ClangDetails) ProtoMessage()               {}
func (*Data_ClangDetails) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{0, 3} }

func (m *Data_ClangDetails) GetPlatforms() int64 {
	if m != nil {
		return m.Platforms
	}
	return 0
}

func (m *Data_ClangDetails) GetCompilationDb() string {
	if m != nil {
		return m.CompilationDb
	}
	return ""
}

func (m *Data_ClangDetails) GetCompDepPath() []string {
	if m != nil {
		return m.CompDepPath
	}
	return nil
}

// Results from running a Tricium analyzer.
//
// Results are returned to the Tricium service via isolated output from
// swarming tasks executing Tricium workers.
//
// ISOLATED PATH: tricium/data/results.json
type Data_Results struct {
	Platforms int64 `protobuf:"varint,1,opt,name=platforms" json:"platforms,omitempty"`
	// Zero or more results found as comments, either inline comments or change
	// comments (comments without line positions).
	Comment []*Data_Comment `protobuf:"bytes,2,rep,name=comment" json:"comment,omitempty"`
}

func (m *Data_Results) Reset()                    { *m = Data_Results{} }
func (m *Data_Results) String() string            { return proto.CompactTextString(m) }
func (*Data_Results) ProtoMessage()               {}
func (*Data_Results) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{0, 4} }

func (m *Data_Results) GetPlatforms() int64 {
	if m != nil {
		return m.Platforms
	}
	return 0
}

func (m *Data_Results) GetComment() []*Data_Comment {
	if m != nil {
		return m.Comment
	}
	return nil
}

// Results.Comment, results as comments.
//
// Similar content as that needed to provide robot comments in Gerrit,
// go/robot-comments-in-gerrit.
type Data_Comment struct {
	// Category of the result, encoded as a path with the analyzer name as the
	// root, annotated with the platform, followed by an arbitrary number of
	// subcategories, e.g., ‘ClangTidy_${platform}/llvm-header-guard’.
	Category string `protobuf:"bytes,1,opt,name=category" json:"category,omitempty"`
	// Comment message.  This should be a short message suitable as a code
	// review comment.
	Message string `protobuf:"bytes,2,opt,name=message" json:"message,omitempty"`
	// URL to more information.  Use this field to keep the message of the
	// comment short.
	Url string `protobuf:"bytes,3,opt,name=url" json:"url,omitempty"`
	// Path to the file this comment is for.
	Path string `protobuf:"bytes,4,opt,name=path" json:"path,omitempty"`
	// Position information. This information can be left out if the comment
	// is on file level.  If position is given, then at least the start line
	// should be included.
	StartLine int32 `protobuf:"varint,5,opt,name=start_line,json=startLine" json:"start_line,omitempty"`
	EndLine   int32 `protobuf:"varint,6,opt,name=end_line,json=endLine" json:"end_line,omitempty"`
	StartChar int32 `protobuf:"varint,7,opt,name=start_char,json=startChar" json:"start_char,omitempty"`
	EndChar   int32 `protobuf:"varint,8,opt,name=end_char,json=endChar" json:"end_char,omitempty"`
	// Suggested fixes for the identified issue.
	Suggestion []*Data_Suggestion `protobuf:"bytes,9,rep,name=suggestion" json:"suggestion,omitempty"`
}

func (m *Data_Comment) Reset()                    { *m = Data_Comment{} }
func (m *Data_Comment) String() string            { return proto.CompactTextString(m) }
func (*Data_Comment) ProtoMessage()               {}
func (*Data_Comment) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{0, 5} }

func (m *Data_Comment) GetCategory() string {
	if m != nil {
		return m.Category
	}
	return ""
}

func (m *Data_Comment) GetMessage() string {
	if m != nil {
		return m.Message
	}
	return ""
}

func (m *Data_Comment) GetUrl() string {
	if m != nil {
		return m.Url
	}
	return ""
}

func (m *Data_Comment) GetPath() string {
	if m != nil {
		return m.Path
	}
	return ""
}

func (m *Data_Comment) GetStartLine() int32 {
	if m != nil {
		return m.StartLine
	}
	return 0
}

func (m *Data_Comment) GetEndLine() int32 {
	if m != nil {
		return m.EndLine
	}
	return 0
}

func (m *Data_Comment) GetStartChar() int32 {
	if m != nil {
		return m.StartChar
	}
	return 0
}

func (m *Data_Comment) GetEndChar() int32 {
	if m != nil {
		return m.EndChar
	}
	return 0
}

func (m *Data_Comment) GetSuggestion() []*Data_Suggestion {
	if m != nil {
		return m.Suggestion
	}
	return nil
}

// Suggested fix.
//
// A fix may include replacements in any file in the same repo as the file of
// the corresponding comment.
type Data_Suggestion struct {
	// A brief description of the suggested fix.
	Description string `protobuf:"bytes,1,opt,name=description" json:"description,omitempty"`
	// Fix as a list of replacements.
	Replacement []*Data_Replacement `protobuf:"bytes,2,rep,name=replacement" json:"replacement,omitempty"`
}

func (m *Data_Suggestion) Reset()                    { *m = Data_Suggestion{} }
func (m *Data_Suggestion) String() string            { return proto.CompactTextString(m) }
func (*Data_Suggestion) ProtoMessage()               {}
func (*Data_Suggestion) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{0, 6} }

func (m *Data_Suggestion) GetDescription() string {
	if m != nil {
		return m.Description
	}
	return ""
}

func (m *Data_Suggestion) GetReplacement() []*Data_Replacement {
	if m != nil {
		return m.Replacement
	}
	return nil
}

// A suggested replacement.
//
// The replacement should be for one continuous section of a file.
type Data_Replacement struct {
	// Path to file for this replacement.
	Path string `protobuf:"bytes,1,opt,name=path" json:"path,omitempty"`
	// A replacement string.
	Replacement string `protobuf:"bytes,2,opt,name=replacement" json:"replacement,omitempty"`
	// A continuous section of the file to replace.
	StartLine int32 `protobuf:"varint,3,opt,name=start_line,json=startLine" json:"start_line,omitempty"`
	EndLine   int32 `protobuf:"varint,4,opt,name=end_line,json=endLine" json:"end_line,omitempty"`
	StartChar int32 `protobuf:"varint,5,opt,name=start_char,json=startChar" json:"start_char,omitempty"`
	EndChar   int32 `protobuf:"varint,6,opt,name=end_char,json=endChar" json:"end_char,omitempty"`
}

func (m *Data_Replacement) Reset()                    { *m = Data_Replacement{} }
func (m *Data_Replacement) String() string            { return proto.CompactTextString(m) }
func (*Data_Replacement) ProtoMessage()               {}
func (*Data_Replacement) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{0, 7} }

func (m *Data_Replacement) GetPath() string {
	if m != nil {
		return m.Path
	}
	return ""
}

func (m *Data_Replacement) GetReplacement() string {
	if m != nil {
		return m.Replacement
	}
	return ""
}

func (m *Data_Replacement) GetStartLine() int32 {
	if m != nil {
		return m.StartLine
	}
	return 0
}

func (m *Data_Replacement) GetEndLine() int32 {
	if m != nil {
		return m.EndLine
	}
	return 0
}

func (m *Data_Replacement) GetStartChar() int32 {
	if m != nil {
		return m.StartChar
	}
	return 0
}

func (m *Data_Replacement) GetEndChar() int32 {
	if m != nil {
		return m.EndChar
	}
	return 0
}

func init() {
	proto.RegisterType((*Data)(nil), "tricium.Data")
	proto.RegisterType((*Data_TypeDetails)(nil), "tricium.Data.TypeDetails")
	proto.RegisterType((*Data_GitFileDetails)(nil), "tricium.Data.GitFileDetails")
	proto.RegisterType((*Data_Files)(nil), "tricium.Data.Files")
	proto.RegisterType((*Data_ClangDetails)(nil), "tricium.Data.ClangDetails")
	proto.RegisterType((*Data_Results)(nil), "tricium.Data.Results")
	proto.RegisterType((*Data_Comment)(nil), "tricium.Data.Comment")
	proto.RegisterType((*Data_Suggestion)(nil), "tricium.Data.Suggestion")
	proto.RegisterType((*Data_Replacement)(nil), "tricium.Data.Replacement")
	proto.RegisterEnum("tricium.Data_Type", Data_Type_name, Data_Type_value)
}

func init() { proto.RegisterFile("infra/tricium/api/v1/data.proto", fileDescriptor1) }

var fileDescriptor1 = []byte{
	// 576 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0x8c, 0x94, 0x4f, 0x6f, 0xd3, 0x30,
	0x18, 0xc6, 0x49, 0x93, 0x2e, 0xcd, 0x1b, 0x36, 0x15, 0x6b, 0x48, 0x59, 0xc4, 0x9f, 0x6a, 0x12,
	0xa8, 0xa7, 0x16, 0xc6, 0x05, 0xc4, 0x69, 0x6a, 0xbb, 0x69, 0x52, 0x35, 0x86, 0x53, 0x24, 0x6e,
	0x91, 0x97, 0xb8, 0x99, 0x45, 0xfe, 0x58, 0xb6, 0x0b, 0xda, 0x99, 0x8f, 0xc0, 0x87, 0xe1, 0xeb,
	0xa1, 0x78, 0x49, 0xe6, 0xee, 0xd0, 0x71, 0x8b, 0x9f, 0xe7, 0x79, 0xfd, 0xbc, 0xfd, 0x55, 0x32,
	0xbc, 0x66, 0xe5, 0x5a, 0x90, 0xa9, 0x12, 0x2c, 0x61, 0x9b, 0x62, 0x4a, 0x38, 0x9b, 0xfe, 0x7c,
	0x3f, 0x4d, 0x89, 0x22, 0x13, 0x2e, 0x2a, 0x55, 0x21, 0xb7, 0xb1, 0x8e, 0x7f, 0x7b, 0xe0, 0xcc,
	0x89, 0x22, 0x61, 0x06, 0xfe, 0xea, 0x96, 0xd3, 0x39, 0x55, 0x84, 0xe5, 0x12, 0xbd, 0x05, 0x47,
	0xdd, 0x72, 0x1a, 0x58, 0x23, 0x6b, 0x7c, 0x70, 0x82, 0x26, 0x4d, 0x7e, 0x52, 0x67, 0x27, 0x75,
	0x10, 0x6b, 0x1f, 0xbd, 0x83, 0x43, 0x26, 0x63, 0x9e, 0x13, 0xb5, 0xae, 0x44, 0x11, 0x4b, 0x4e,
	0x13, 0xb6, 0x66, 0x49, 0xd0, 0x1b, 0x59, 0xe3, 0x01, 0x46, 0x4c, 0x5e, 0x35, 0x56, 0xd4, 0x38,
	0xa1, 0x82, 0x83, 0x73, 0xa6, 0xce, 0x58, 0xde, 0x75, 0xbd, 0x00, 0xaf, 0xbd, 0x40, 0xea, 0x42,
	0x1b, 0xdf, 0x0b, 0xe8, 0x15, 0x80, 0xa0, 0xbc, 0x92, 0x4c, 0x55, 0xe2, 0x56, 0xdf, 0xeb, 0x61,
	0x43, 0x41, 0x43, 0xb0, 0x05, 0x5d, 0x07, 0xb6, 0x36, 0xea, 0x4f, 0x84, 0xc0, 0xe1, 0x44, 0xdd,
	0x04, 0xce, 0xc8, 0x1e, 0x7b, 0x58, 0x7f, 0x87, 0x9f, 0xa0, 0x5f, 0x57, 0x3e, 0x56, 0xd6, 0x8e,
	0xf6, 0x8c, 0xd1, 0x5f, 0xf0, 0x74, 0x96, 0x93, 0x32, 0xfb, 0xbf, 0x75, 0xdf, 0xc0, 0x41, 0x52,
	0x15, 0x9c, 0xe5, 0x44, 0xb1, 0xaa, 0x8c, 0xd3, 0xeb, 0x66, 0xe5, 0x7d, 0x43, 0x9d, 0x5f, 0xa3,
	0x63, 0xd0, 0x42, 0x9c, 0x52, 0x1e, 0xeb, 0x46, 0x5b, 0x37, 0xfa, 0xb5, 0x38, 0xa7, 0xfc, 0xaa,
	0x2e, 0xfe, 0x0e, 0x2e, 0xa6, 0x72, 0x93, 0xab, 0xc7, 0x3a, 0xa7, 0xe0, 0x26, 0x55, 0x51, 0xd0,
	0x52, 0xe9, 0xc5, 0xfd, 0x93, 0xe7, 0xdb, 0xff, 0xd7, 0xec, 0xce, 0xc4, 0x6d, 0x2a, 0xfc, 0xd3,
	0x03, 0xb7, 0x11, 0x51, 0x08, 0x83, 0x84, 0x28, 0x9a, 0xd5, 0x74, 0x2d, 0xbd, 0x6a, 0x77, 0x46,
	0x01, 0xb8, 0x05, 0x95, 0x92, 0x64, 0xb4, 0xf9, 0x15, 0xed, 0xb1, 0xa6, 0xbe, 0x11, 0x79, 0x4b,
	0x7d, 0x23, 0x72, 0x83, 0xba, 0xd5, 0xa2, 0x43, 0x2f, 0x01, 0xa4, 0x22, 0x42, 0xc5, 0x39, 0x2b,
	0x69, 0xd0, 0x1f, 0x59, 0xe3, 0x3e, 0xf6, 0xb4, 0xb2, 0x64, 0x25, 0x45, 0x47, 0x30, 0xa0, 0x65,
	0x7a, 0x67, 0xee, 0x69, 0xd3, 0xa5, 0x65, 0xaa, 0xad, 0x6e, 0x32, 0xb9, 0x21, 0x22, 0x70, 0x8d,
	0xc9, 0xd9, 0x0d, 0x11, 0xed, 0xa4, 0x36, 0x07, 0xdd, 0xa4, 0xb6, 0x3e, 0x02, 0xc8, 0x4d, 0x96,
	0x51, 0x59, 0x93, 0x0e, 0x3c, 0xcd, 0x23, 0xd8, 0xe6, 0x11, 0x75, 0x3e, 0x36, 0xb2, 0xe1, 0x0f,
	0x80, 0x7b, 0x07, 0x8d, 0xc0, 0x4f, 0xa9, 0x4c, 0x04, 0xe3, 0xfa, 0xa2, 0x3b, 0x34, 0xa6, 0x84,
	0x3e, 0x83, 0x2f, 0x28, 0xcf, 0x49, 0x42, 0x0d, 0xf4, 0x47, 0xdb, 0x55, 0xf8, 0x3e, 0x80, 0xcd,
	0x74, 0xf8, 0xd7, 0x02, 0xdf, 0x30, 0x3b, 0x7c, 0x96, 0x81, 0x6f, 0xf4, 0xb0, 0x40, 0xaf, 0x60,
	0x48, 0x0f, 0x00, 0xdb, 0xbb, 0x00, 0x3b, 0xbb, 0x00, 0xf7, 0x77, 0x01, 0xde, 0xdb, 0x02, 0x7c,
	0xfc, 0x15, 0x9c, 0xfa, 0x01, 0x40, 0x03, 0x70, 0x2e, 0xbf, 0x5c, 0x2e, 0x86, 0x4f, 0xd0, 0x21,
	0x0c, 0xcf, 0x2f, 0x56, 0xf1, 0xd9, 0xc5, 0x72, 0x11, 0xcf, 0x17, 0xab, 0xd3, 0x8b, 0x65, 0x34,
	0xb4, 0x90, 0x07, 0xfd, 0x5a, 0x89, 0x86, 0x3d, 0xf4, 0x0c, 0xf6, 0x67, 0xcb, 0xd3, 0xcb, 0xf3,
	0xce, 0xb5, 0x91, 0x0f, 0x2e, 0x5e, 0x44, 0xdf, 0x96, 0xab, 0x68, 0xe8, 0x5c, 0xef, 0xe9, 0x57,
	0xe9, 0xc3, 0xbf, 0x00, 0x00, 0x00, 0xff, 0xff, 0xf9, 0xf6, 0xec, 0x48, 0xb8, 0x04, 0x00, 0x00,
}
