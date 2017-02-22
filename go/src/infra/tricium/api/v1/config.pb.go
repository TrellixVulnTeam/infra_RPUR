// Code generated by protoc-gen-go.
// source: infra/tricium/api/v1/config.proto
// DO NOT EDIT!

/*
Package tricium is a generated protocol buffer package.

It is generated from these files:
	infra/tricium/api/v1/config.proto
	infra/tricium/api/v1/data.proto
	infra/tricium/api/v1/platform.proto
	infra/tricium/api/v1/tricium.proto

It has these top-level messages:
	ServiceConfig
	ProjectDetails
	ProjectConfig
	RepoDetails
	GitRepoDetails
	Acl
	Selection
	Analyzer
	ConfigDef
	Impl
	Recipe
	Property
	Config
	Cmd
	CipdPackage
	Data
	Platform
	AnalyzeRequest
	AnalyzeResponse
*/
package tricium

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

// Supported kinds of repositories.
type RepoDetails_Kind int32

const (
	RepoDetails_GIT RepoDetails_Kind = 0
)

var RepoDetails_Kind_name = map[int32]string{
	0: "GIT",
}
var RepoDetails_Kind_value = map[string]int32{
	"GIT": 0,
}

func (x RepoDetails_Kind) String() string {
	return proto.EnumName(RepoDetails_Kind_name, int32(x))
}
func (RepoDetails_Kind) EnumDescriptor() ([]byte, []int) { return fileDescriptor0, []int{3, 0} }

// Roles relevant to Tricium.
type Acl_Role int32

const (
	// Can read progress/results.
	Acl_READER Acl_Role = 0
	// Can request analysis.
	Acl_REQUESTER Acl_Role = 1
)

var Acl_Role_name = map[int32]string{
	0: "READER",
	1: "REQUESTER",
}
var Acl_Role_value = map[string]int32{
	"READER":    0,
	"REQUESTER": 1,
}

func (x Acl_Role) String() string {
	return proto.EnumName(Acl_Role_name, int32(x))
}
func (Acl_Role) EnumDescriptor() ([]byte, []int) { return fileDescriptor0, []int{5, 0} }

// Tricium service configuration.
//
// Listing supported platforms and analyzers shared between projects connected
// to Tricium.
type ServiceConfig struct {
	// Supported platforms.
	Platforms []*Platform_Details `protobuf:"bytes,1,rep,name=platforms" json:"platforms,omitempty"`
	// Supported data types.
	DataDetails []*Data_TypeDetails `protobuf:"bytes,2,rep,name=data_details,json=dataDetails" json:"data_details,omitempty"`
	// List of shared analyzers.
	Analyzers []*Analyzer `protobuf:"bytes,3,rep,name=analyzers" json:"analyzers,omitempty"`
	// Details for connected projects.
	Projects []*ProjectDetails `protobuf:"bytes,4,rep,name=projects" json:"projects,omitempty"`
	// PubSub topic used to collect worker completion notifications from Swarming.
	SwarmingWorkerTopic string `protobuf:"bytes,5,opt,name=swarming_worker_topic,json=swarmingWorkerTopic" json:"swarming_worker_topic,omitempty"`
	// Base recipe command used for workers implemented as recipes.
	//
	// Specific recipe details for the worker will be added as flags at the
	// end of the argument list.
	RecipeCmd *Cmd `protobuf:"bytes,6,opt,name=RecipeCmd,json=recipeCmd" json:"RecipeCmd,omitempty"`
	// Base recipe packages used for workers implemented as recipes.
	//
	// These packages will be adjusted for the platform in question, by appending
	// platform name details to the end of the package name.
	RecipePackages []*CipdPackage `protobuf:"bytes,7,rep,name=RecipePackages,json=recipePackages" json:"RecipePackages,omitempty"`
}

func (m *ServiceConfig) Reset()                    { *m = ServiceConfig{} }
func (m *ServiceConfig) String() string            { return proto.CompactTextString(m) }
func (*ServiceConfig) ProtoMessage()               {}
func (*ServiceConfig) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *ServiceConfig) GetPlatforms() []*Platform_Details {
	if m != nil {
		return m.Platforms
	}
	return nil
}

func (m *ServiceConfig) GetDataDetails() []*Data_TypeDetails {
	if m != nil {
		return m.DataDetails
	}
	return nil
}

func (m *ServiceConfig) GetAnalyzers() []*Analyzer {
	if m != nil {
		return m.Analyzers
	}
	return nil
}

func (m *ServiceConfig) GetProjects() []*ProjectDetails {
	if m != nil {
		return m.Projects
	}
	return nil
}

func (m *ServiceConfig) GetSwarmingWorkerTopic() string {
	if m != nil {
		return m.SwarmingWorkerTopic
	}
	return ""
}

func (m *ServiceConfig) GetRecipeCmd() *Cmd {
	if m != nil {
		return m.RecipeCmd
	}
	return nil
}

func (m *ServiceConfig) GetRecipePackages() []*CipdPackage {
	if m != nil {
		return m.RecipePackages
	}
	return nil
}

type ProjectDetails struct {
	// Project name used to map these project details to the config for a project.
	Name string `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
	// General service account for this project.
	// Used for any service interaction, with the exception of swarming.
	ServiceAccount string `protobuf:"bytes,2,opt,name=service_account,json=serviceAccount" json:"service_account,omitempty"`
	// Project-specific swarming service account.
	SwarmingServiceAccount string `protobuf:"bytes,3,opt,name=swarming_service_account,json=swarmingServiceAccount" json:"swarming_service_account,omitempty"`
	// Details of the repository connected to the project. This should be the
	// repository hosting the files that should be analyzed for this project.
	RepoDetails *RepoDetails `protobuf:"bytes,4,opt,name=repo_details,json=repoDetails" json:"repo_details,omitempty"`
}

func (m *ProjectDetails) Reset()                    { *m = ProjectDetails{} }
func (m *ProjectDetails) String() string            { return proto.CompactTextString(m) }
func (*ProjectDetails) ProtoMessage()               {}
func (*ProjectDetails) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *ProjectDetails) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *ProjectDetails) GetServiceAccount() string {
	if m != nil {
		return m.ServiceAccount
	}
	return ""
}

func (m *ProjectDetails) GetSwarmingServiceAccount() string {
	if m != nil {
		return m.SwarmingServiceAccount
	}
	return ""
}

func (m *ProjectDetails) GetRepoDetails() *RepoDetails {
	if m != nil {
		return m.RepoDetails
	}
	return nil
}

// Tricium project configuration.
//
// Specifies details needed to connect a project to Tricium, adds project
// specific analyzers and implementations, and selects analyzer
// implementations.
type ProjectConfig struct {
	// Project name,
	Name string `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
	// Access control rules for the project.
	Acls []*Acl `protobuf:"bytes,2,rep,name=acls" json:"acls,omitempty"`
	// Project-specific analyzer details. This includes project-specific analyzer
	// implementations and full project-specific analyzer specifications.
	Analyzers []*Analyzer `protobuf:"bytes,3,rep,name=analyzers" json:"analyzers,omitempty"`
	// Selection of analyzer implementations to run for this project.
	Selections []*Selection `protobuf:"bytes,4,rep,name=selections" json:"selections,omitempty"`
}

func (m *ProjectConfig) Reset()                    { *m = ProjectConfig{} }
func (m *ProjectConfig) String() string            { return proto.CompactTextString(m) }
func (*ProjectConfig) ProtoMessage()               {}
func (*ProjectConfig) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *ProjectConfig) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *ProjectConfig) GetAcls() []*Acl {
	if m != nil {
		return m.Acls
	}
	return nil
}

func (m *ProjectConfig) GetAnalyzers() []*Analyzer {
	if m != nil {
		return m.Analyzers
	}
	return nil
}

func (m *ProjectConfig) GetSelections() []*Selection {
	if m != nil {
		return m.Selections
	}
	return nil
}

// Repository details for a project.
type RepoDetails struct {
	Kind RepoDetails_Kind `protobuf:"varint,1,opt,name=kind,enum=tricium.RepoDetails_Kind" json:"kind,omitempty"`
	// If repository kind is GIT then provide Git details.
	GitDetails *GitRepoDetails `protobuf:"bytes,2,opt,name=git_details,json=gitDetails" json:"git_details,omitempty"`
}

func (m *RepoDetails) Reset()                    { *m = RepoDetails{} }
func (m *RepoDetails) String() string            { return proto.CompactTextString(m) }
func (*RepoDetails) ProtoMessage()               {}
func (*RepoDetails) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func (m *RepoDetails) GetKind() RepoDetails_Kind {
	if m != nil {
		return m.Kind
	}
	return RepoDetails_GIT
}

func (m *RepoDetails) GetGitDetails() *GitRepoDetails {
	if m != nil {
		return m.GitDetails
	}
	return nil
}

// Git repository details.
type GitRepoDetails struct {
	// URL to repository.
	Repository string `protobuf:"bytes,1,opt,name=repository" json:"repository,omitempty"`
	// Default ref to use to get files to analyze.
	Ref string `protobuf:"bytes,2,opt,name=ref" json:"ref,omitempty"`
}

func (m *GitRepoDetails) Reset()                    { *m = GitRepoDetails{} }
func (m *GitRepoDetails) String() string            { return proto.CompactTextString(m) }
func (*GitRepoDetails) ProtoMessage()               {}
func (*GitRepoDetails) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

func (m *GitRepoDetails) GetRepository() string {
	if m != nil {
		return m.Repository
	}
	return ""
}

func (m *GitRepoDetails) GetRef() string {
	if m != nil {
		return m.Ref
	}
	return ""
}

// Access control rules.
type Acl struct {
	// Role of a group or identity.
	Role Acl_Role `protobuf:"varint,1,opt,name=role,enum=tricium.Acl_Role" json:"role,omitempty"`
	// Name of group, as defined in the auth service. Specify either group or
	// identity, not both.
	Group string `protobuf:"bytes,2,opt,name=group" json:"group,omitempty"`
	// Identity, as defined by the auth service. Can be either an email address
	// or an indentity string, for instance, "anonymous:anonymous" for anonymous
	// users. Specify either group or identity, not both.
	Identity string `protobuf:"bytes,3,opt,name=identity" json:"identity,omitempty"`
}

func (m *Acl) Reset()                    { *m = Acl{} }
func (m *Acl) String() string            { return proto.CompactTextString(m) }
func (*Acl) ProtoMessage()               {}
func (*Acl) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{5} }

func (m *Acl) GetRole() Acl_Role {
	if m != nil {
		return m.Role
	}
	return Acl_READER
}

func (m *Acl) GetGroup() string {
	if m != nil {
		return m.Group
	}
	return ""
}

func (m *Acl) GetIdentity() string {
	if m != nil {
		return m.Identity
	}
	return ""
}

// Selection of analyzer implementations to run for a project.
type Selection struct {
	// Name of analyzer to run.
	Analyzer string `protobuf:"bytes,1,opt,name=analyzer" json:"analyzer,omitempty"`
	// Name of platform to retrieve results from.
	Platform Platform_Name `protobuf:"varint,2,opt,name=platform,enum=tricium.Platform_Name" json:"platform,omitempty"`
	// Analyzer configuration to use on this platform.
	Configs []*Config `protobuf:"bytes,3,rep,name=configs" json:"configs,omitempty"`
}

func (m *Selection) Reset()                    { *m = Selection{} }
func (m *Selection) String() string            { return proto.CompactTextString(m) }
func (*Selection) ProtoMessage()               {}
func (*Selection) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{6} }

func (m *Selection) GetAnalyzer() string {
	if m != nil {
		return m.Analyzer
	}
	return ""
}

func (m *Selection) GetPlatform() Platform_Name {
	if m != nil {
		return m.Platform
	}
	return Platform_ANY
}

func (m *Selection) GetConfigs() []*Config {
	if m != nil {
		return m.Configs
	}
	return nil
}

// Analyzer specification.
type Analyzer struct {
	// Name of analyzer. This name is used to select the analyzer and is used
	// when reporting results for the analyzer. This name should be unique among
	// Tricium analyzers.
	Name string `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
	// Tricium data needed by this analyzer.
	Needs Data_Type `protobuf:"varint,2,opt,name=needs,enum=tricium.Data_Type" json:"needs,omitempty"`
	// Tricium data provided by this analyzer.
	Provides Data_Type `protobuf:"varint,3,opt,name=provides,enum=tricium.Data_Type" json:"provides,omitempty"`
	// Paths to run this analyzer on, defined as a glob.
	PathFilters []string `protobuf:"bytes,4,rep,name=path_filters,json=pathFilters" json:"path_filters,omitempty"`
	// Email to the owner of this analyzer.
	Owner string `protobuf:"bytes,6,opt,name=owner" json:"owner,omitempty"`
	// Monorail bug component for bug filing.
	Component string `protobuf:"bytes,7,opt,name=component" json:"component,omitempty"`
	// Analyzer configuration. These configuration options enable projects to
	// customize how an analyzer implementation analyzes their files.  It's
	// common for analyzers to provide a list of possible checks which can be
	// configured via a command line flag or similar. This field provides a way
	// to expose such flags as configuration options.
	ConfigDefs []*ConfigDef `protobuf:"bytes,8,rep,name=config_defs,json=configDefs" json:"config_defs,omitempty"`
	// Analyzer implementations. An analyzer may run on many platforms and this
	// may require many different implementations of the analyzer. An
	// implementation may be shared between several platforms if possible.
	Impls []*Impl `protobuf:"bytes,9,rep,name=impls" json:"impls,omitempty"`
}

func (m *Analyzer) Reset()                    { *m = Analyzer{} }
func (m *Analyzer) String() string            { return proto.CompactTextString(m) }
func (*Analyzer) ProtoMessage()               {}
func (*Analyzer) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{7} }

func (m *Analyzer) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *Analyzer) GetNeeds() Data_Type {
	if m != nil {
		return m.Needs
	}
	return Data_NONE
}

func (m *Analyzer) GetProvides() Data_Type {
	if m != nil {
		return m.Provides
	}
	return Data_NONE
}

func (m *Analyzer) GetPathFilters() []string {
	if m != nil {
		return m.PathFilters
	}
	return nil
}

func (m *Analyzer) GetOwner() string {
	if m != nil {
		return m.Owner
	}
	return ""
}

func (m *Analyzer) GetComponent() string {
	if m != nil {
		return m.Component
	}
	return ""
}

func (m *Analyzer) GetConfigDefs() []*ConfigDef {
	if m != nil {
		return m.ConfigDefs
	}
	return nil
}

func (m *Analyzer) GetImpls() []*Impl {
	if m != nil {
		return m.Impls
	}
	return nil
}

// Definition of an analyzer configuration, e.g., ClangTidy is configured with
// a ‘checks’ flag.
type ConfigDef struct {
	// Name of configuration option.
	Name string `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
	// Default value for the config, e.g., checks=”all”.
	Default string `protobuf:"bytes,2,opt,name=default" json:"default,omitempty"`
}

func (m *ConfigDef) Reset()                    { *m = ConfigDef{} }
func (m *ConfigDef) String() string            { return proto.CompactTextString(m) }
func (*ConfigDef) ProtoMessage()               {}
func (*ConfigDef) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{8} }

func (m *ConfigDef) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *ConfigDef) GetDefault() string {
	if m != nil {
		return m.Default
	}
	return ""
}

// Analyzer implementation.
//
// Implementation can be either recipe-based or command-based.
// If platform-specific data is needed or provided, the specific platform
// details should be provided in the implementation.
// Note that the runtime platform of the implementation may be different than
// the platform(s) used to refine the data-dependency.
type Impl struct {
	// Data-dependency details specific to this implementation.
	// For instance, if the needed data needs to be tied to a specific platform
	// then the 'needs_for_platform' field should be set to that platform.
	// Likewise for any provided data type tied to a specific platform, this
	// should be indicated with the 'provides_for_platform' field.
	// Either if these fields can be left out for implementations of analyzers
	// not needing or providing platform-specific data.
	NeedsForPlatform    Platform_Name `protobuf:"varint,1,opt,name=needs_for_platform,json=needsForPlatform,enum=tricium.Platform_Name" json:"needs_for_platform,omitempty"`
	ProvidesForPlatform Platform_Name `protobuf:"varint,2,opt,name=provides_for_platform,json=providesForPlatform,enum=tricium.Platform_Name" json:"provides_for_platform,omitempty"`
	// The platform to run this implementation on. This may be different
	// from the platforms used to refine data-dependencies, as long as the
	// data consumed/produced follows the specification.
	RuntimePlatform Platform_Name `protobuf:"varint,3,opt,name=runtime_platform,json=runtimePlatform,enum=tricium.Platform_Name" json:"runtime_platform,omitempty"`
	// Cipd packages needed by this implementation.
	CipdPackages []*CipdPackage `protobuf:"bytes,4,rep,name=cipd_packages,json=cipdPackages" json:"cipd_packages,omitempty"`
	// Types that are valid to be assigned to Impl:
	//	*Impl_Recipe
	//	*Impl_Cmd
	Impl isImpl_Impl `protobuf_oneof:"impl"`
	// Deadline for execution of corresponding worker (in minutes). Note that
	// this deadline includes the launch of a swarming task for the corresponding
	// worker, and collection of results from that worker.
	Deadline int32 `protobuf:"varint,7,opt,name=deadline" json:"deadline,omitempty"`
}

func (m *Impl) Reset()                    { *m = Impl{} }
func (m *Impl) String() string            { return proto.CompactTextString(m) }
func (*Impl) ProtoMessage()               {}
func (*Impl) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{9} }

type isImpl_Impl interface {
	isImpl_Impl()
}

type Impl_Recipe struct {
	Recipe *Recipe `protobuf:"bytes,5,opt,name=recipe,oneof"`
}
type Impl_Cmd struct {
	Cmd *Cmd `protobuf:"bytes,6,opt,name=cmd,oneof"`
}

func (*Impl_Recipe) isImpl_Impl() {}
func (*Impl_Cmd) isImpl_Impl()    {}

func (m *Impl) GetImpl() isImpl_Impl {
	if m != nil {
		return m.Impl
	}
	return nil
}

func (m *Impl) GetNeedsForPlatform() Platform_Name {
	if m != nil {
		return m.NeedsForPlatform
	}
	return Platform_ANY
}

func (m *Impl) GetProvidesForPlatform() Platform_Name {
	if m != nil {
		return m.ProvidesForPlatform
	}
	return Platform_ANY
}

func (m *Impl) GetRuntimePlatform() Platform_Name {
	if m != nil {
		return m.RuntimePlatform
	}
	return Platform_ANY
}

func (m *Impl) GetCipdPackages() []*CipdPackage {
	if m != nil {
		return m.CipdPackages
	}
	return nil
}

func (m *Impl) GetRecipe() *Recipe {
	if x, ok := m.GetImpl().(*Impl_Recipe); ok {
		return x.Recipe
	}
	return nil
}

func (m *Impl) GetCmd() *Cmd {
	if x, ok := m.GetImpl().(*Impl_Cmd); ok {
		return x.Cmd
	}
	return nil
}

func (m *Impl) GetDeadline() int32 {
	if m != nil {
		return m.Deadline
	}
	return 0
}

// XXX_OneofFuncs is for the internal use of the proto package.
func (*Impl) XXX_OneofFuncs() (func(msg proto.Message, b *proto.Buffer) error, func(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error), func(msg proto.Message) (n int), []interface{}) {
	return _Impl_OneofMarshaler, _Impl_OneofUnmarshaler, _Impl_OneofSizer, []interface{}{
		(*Impl_Recipe)(nil),
		(*Impl_Cmd)(nil),
	}
}

func _Impl_OneofMarshaler(msg proto.Message, b *proto.Buffer) error {
	m := msg.(*Impl)
	// impl
	switch x := m.Impl.(type) {
	case *Impl_Recipe:
		b.EncodeVarint(5<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.Recipe); err != nil {
			return err
		}
	case *Impl_Cmd:
		b.EncodeVarint(6<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.Cmd); err != nil {
			return err
		}
	case nil:
	default:
		return fmt.Errorf("Impl.Impl has unexpected type %T", x)
	}
	return nil
}

func _Impl_OneofUnmarshaler(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error) {
	m := msg.(*Impl)
	switch tag {
	case 5: // impl.recipe
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(Recipe)
		err := b.DecodeMessage(msg)
		m.Impl = &Impl_Recipe{msg}
		return true, err
	case 6: // impl.cmd
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(Cmd)
		err := b.DecodeMessage(msg)
		m.Impl = &Impl_Cmd{msg}
		return true, err
	default:
		return false, nil
	}
}

func _Impl_OneofSizer(msg proto.Message) (n int) {
	m := msg.(*Impl)
	// impl
	switch x := m.Impl.(type) {
	case *Impl_Recipe:
		s := proto.Size(x.Recipe)
		n += proto.SizeVarint(5<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case *Impl_Cmd:
		s := proto.Size(x.Cmd)
		n += proto.SizeVarint(6<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case nil:
	default:
		panic(fmt.Sprintf("proto: unexpected type %T in oneof", x))
	}
	return n
}

// Specification of how to find a recipe.
type Recipe struct {
	// Repository URL of the recipe package.
	Repository string `protobuf:"bytes,1,opt,name=repository" json:"repository,omitempty"`
	// Path to recipe in the repository.
	Path string `protobuf:"bytes,2,opt,name=path" json:"path,omitempty"`
	// Revision to use.
	Revision string `protobuf:"bytes,3,opt,name=revision" json:"revision,omitempty"`
	// Recipe properties which will be provided as a JSON string to the recipe.
	Properties []*Property `protobuf:"bytes,4,rep,name=properties" json:"properties,omitempty"`
}

func (m *Recipe) Reset()                    { *m = Recipe{} }
func (m *Recipe) String() string            { return proto.CompactTextString(m) }
func (*Recipe) ProtoMessage()               {}
func (*Recipe) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{10} }

func (m *Recipe) GetRepository() string {
	if m != nil {
		return m.Repository
	}
	return ""
}

func (m *Recipe) GetPath() string {
	if m != nil {
		return m.Path
	}
	return ""
}

func (m *Recipe) GetRevision() string {
	if m != nil {
		return m.Revision
	}
	return ""
}

func (m *Recipe) GetProperties() []*Property {
	if m != nil {
		return m.Properties
	}
	return nil
}

// Property used to configure a recipe.
//
// All properties for a recipe are encoded as a JSON string and passed via
// the 'kitchen cook' command, ending up as build properties.
type Property struct {
	Key   string `protobuf:"bytes,1,opt,name=key" json:"key,omitempty"`
	Value string `protobuf:"bytes,2,opt,name=value" json:"value,omitempty"`
}

func (m *Property) Reset()                    { *m = Property{} }
func (m *Property) String() string            { return proto.CompactTextString(m) }
func (*Property) ProtoMessage()               {}
func (*Property) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{11} }

func (m *Property) GetKey() string {
	if m != nil {
		return m.Key
	}
	return ""
}

func (m *Property) GetValue() string {
	if m != nil {
		return m.Value
	}
	return ""
}

// Analyzer configuration used when selecting an analyzer implementation.
type Config struct {
	// Name of the configuration option.
	Name string `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
	// Value of the configuration.
	Value string `protobuf:"bytes,2,opt,name=value" json:"value,omitempty"`
}

func (m *Config) Reset()                    { *m = Config{} }
func (m *Config) String() string            { return proto.CompactTextString(m) }
func (*Config) ProtoMessage()               {}
func (*Config) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{12} }

func (m *Config) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *Config) GetValue() string {
	if m != nil {
		return m.Value
	}
	return ""
}

// Specification of a command.
type Cmd struct {
	// Executable binary.
	Exec string `protobuf:"bytes,1,opt,name=exec" json:"exec,omitempty"`
	// Arguments in order.
	Args []string `protobuf:"bytes,2,rep,name=args" json:"args,omitempty"`
}

func (m *Cmd) Reset()                    { *m = Cmd{} }
func (m *Cmd) String() string            { return proto.CompactTextString(m) }
func (*Cmd) ProtoMessage()               {}
func (*Cmd) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{13} }

func (m *Cmd) GetExec() string {
	if m != nil {
		return m.Exec
	}
	return ""
}

func (m *Cmd) GetArgs() []string {
	if m != nil {
		return m.Args
	}
	return nil
}

// CIPD package.
type CipdPackage struct {
	// CIPD package name.
	PackageName string `protobuf:"bytes,1,opt,name=package_name,json=packageName" json:"package_name,omitempty"`
	// Path to directory, relative to the working directory, where to install
	// package. Cannot be empty or start with a slash.
	Path string `protobuf:"bytes,2,opt,name=path" json:"path,omitempty"`
	// Package version.
	Version string `protobuf:"bytes,3,opt,name=version" json:"version,omitempty"`
}

func (m *CipdPackage) Reset()                    { *m = CipdPackage{} }
func (m *CipdPackage) String() string            { return proto.CompactTextString(m) }
func (*CipdPackage) ProtoMessage()               {}
func (*CipdPackage) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{14} }

func (m *CipdPackage) GetPackageName() string {
	if m != nil {
		return m.PackageName
	}
	return ""
}

func (m *CipdPackage) GetPath() string {
	if m != nil {
		return m.Path
	}
	return ""
}

func (m *CipdPackage) GetVersion() string {
	if m != nil {
		return m.Version
	}
	return ""
}

func init() {
	proto.RegisterType((*ServiceConfig)(nil), "tricium.ServiceConfig")
	proto.RegisterType((*ProjectDetails)(nil), "tricium.ProjectDetails")
	proto.RegisterType((*ProjectConfig)(nil), "tricium.ProjectConfig")
	proto.RegisterType((*RepoDetails)(nil), "tricium.RepoDetails")
	proto.RegisterType((*GitRepoDetails)(nil), "tricium.GitRepoDetails")
	proto.RegisterType((*Acl)(nil), "tricium.Acl")
	proto.RegisterType((*Selection)(nil), "tricium.Selection")
	proto.RegisterType((*Analyzer)(nil), "tricium.Analyzer")
	proto.RegisterType((*ConfigDef)(nil), "tricium.ConfigDef")
	proto.RegisterType((*Impl)(nil), "tricium.Impl")
	proto.RegisterType((*Recipe)(nil), "tricium.Recipe")
	proto.RegisterType((*Property)(nil), "tricium.Property")
	proto.RegisterType((*Config)(nil), "tricium.Config")
	proto.RegisterType((*Cmd)(nil), "tricium.Cmd")
	proto.RegisterType((*CipdPackage)(nil), "tricium.CipdPackage")
	proto.RegisterEnum("tricium.RepoDetails_Kind", RepoDetails_Kind_name, RepoDetails_Kind_value)
	proto.RegisterEnum("tricium.Acl_Role", Acl_Role_name, Acl_Role_value)
}

func init() { proto.RegisterFile("infra/tricium/api/v1/config.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 1041 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0x94, 0x56, 0x5f, 0x6f, 0xdb, 0x36,
	0x10, 0x8f, 0x2a, 0xc5, 0xb6, 0xce, 0x89, 0xe3, 0xb0, 0x7f, 0xa6, 0x05, 0xc3, 0xe6, 0xa8, 0x18,
	0x96, 0x0e, 0xa8, 0x8d, 0x3a, 0x0f, 0x6d, 0x81, 0xbe, 0xb8, 0xb1, 0xfb, 0x67, 0x03, 0x86, 0x8e,
	0xc9, 0xb0, 0xb7, 0x09, 0x9c, 0x44, 0x7b, 0x5c, 0x24, 0x51, 0xa0, 0x64, 0x67, 0xde, 0xdb, 0x5e,
	0x56, 0x60, 0xfb, 0x24, 0xdb, 0xa7, 0xd8, 0x47, 0x1b, 0x48, 0x8a, 0x92, 0x9c, 0x1a, 0x1e, 0xf6,
	0xc6, 0xbb, 0xdf, 0xef, 0x8e, 0xe4, 0xdd, 0x4f, 0x47, 0xc1, 0x29, 0x4b, 0xe7, 0x82, 0x8c, 0x0a,
	0xc1, 0x42, 0xb6, 0x4c, 0x46, 0x24, 0x63, 0xa3, 0xd5, 0x93, 0x51, 0xc8, 0xd3, 0x39, 0x5b, 0x0c,
	0x33, 0xc1, 0x0b, 0x8e, 0xda, 0x25, 0x78, 0xf2, 0xd9, 0x56, 0x6e, 0x44, 0x0a, 0xa2, 0x99, 0x27,
	0x0f, 0xb7, 0x12, 0xb2, 0x98, 0x14, 0x73, 0x2e, 0x12, 0x4d, 0xf2, 0xff, 0xb0, 0xe1, 0xf0, 0x92,
	0x8a, 0x15, 0x0b, 0xe9, 0x85, 0xda, 0x06, 0x3d, 0x05, 0xd7, 0x70, 0x72, 0xcf, 0x1a, 0xd8, 0x67,
	0xdd, 0xf1, 0xc7, 0xc3, 0x32, 0xc9, 0xf0, 0x9d, 0x89, 0x9e, 0xd2, 0x82, 0xb0, 0x38, 0xc7, 0x35,
	0x17, 0xbd, 0x80, 0x03, 0xb9, 0x7b, 0x10, 0x69, 0xc8, 0xbb, 0x73, 0x2b, 0x76, 0x2a, 0x8f, 0x76,
	0xb5, 0xce, 0xa8, 0x89, 0xed, 0x4a, 0x7a, 0x69, 0xa0, 0x11, 0xb8, 0x24, 0x25, 0xf1, 0xfa, 0x57,
	0x2a, 0x72, 0xcf, 0x56, 0xa1, 0xc7, 0x55, 0xe8, 0xa4, 0x44, 0x70, 0xcd, 0x41, 0xe7, 0xd0, 0xc9,
	0x04, 0xff, 0x99, 0x86, 0x45, 0xee, 0x39, 0x8a, 0xff, 0x51, 0x7d, 0x4c, 0x0d, 0x98, 0x8d, 0x2a,
	0x22, 0x1a, 0xc3, 0xfd, 0xfc, 0x86, 0x88, 0x84, 0xa5, 0x8b, 0xe0, 0x86, 0x8b, 0x6b, 0x2a, 0x82,
	0x82, 0x67, 0x2c, 0xf4, 0xf6, 0x07, 0xd6, 0x99, 0x8b, 0xef, 0x1a, 0xf0, 0x7b, 0x85, 0x5d, 0x49,
	0x08, 0x7d, 0x09, 0x2e, 0xa6, 0x21, 0xcb, 0xe8, 0x45, 0x12, 0x79, 0xad, 0x81, 0x75, 0xd6, 0x1d,
	0x1f, 0x54, 0x3b, 0x5d, 0x24, 0x11, 0x76, 0x85, 0x81, 0xd1, 0x0b, 0xe8, 0x69, 0xee, 0x3b, 0x12,
	0x5e, 0x93, 0x05, 0xcd, 0xbd, 0xb6, 0x3a, 0xda, 0xbd, 0x3a, 0x80, 0x65, 0x51, 0x09, 0xe2, 0x9e,
	0xd8, 0xe0, 0xfa, 0xff, 0x58, 0xd0, 0xdb, 0x3c, 0x3a, 0x42, 0xe0, 0xa4, 0x24, 0xa1, 0x9e, 0xa5,
	0xce, 0xa7, 0xd6, 0xe8, 0x0b, 0x38, 0xca, 0x75, 0xcb, 0x02, 0x12, 0x86, 0x7c, 0x99, 0x16, 0xde,
	0x1d, 0x05, 0xf7, 0x4a, 0xf7, 0x44, 0x7b, 0xd1, 0x33, 0xf0, 0xaa, 0xdb, 0xde, 0x8e, 0xb0, 0x55,
	0xc4, 0x03, 0x83, 0x5f, 0x6e, 0x46, 0x3e, 0x85, 0x03, 0x41, 0x33, 0x5e, 0xf5, 0xd2, 0x51, 0xd7,
	0xae, 0x6f, 0x81, 0x69, 0xc6, 0xab, 0x36, 0x8a, 0xda, 0xf0, 0xff, 0xb6, 0xe0, 0xb0, 0xbc, 0x42,
	0xa9, 0xa7, 0x6d, 0x37, 0x18, 0x80, 0x43, 0xc2, 0x4a, 0x22, 0x75, 0x35, 0x27, 0x61, 0x8c, 0x15,
	0xf2, 0xff, 0xe5, 0x30, 0x06, 0xc8, 0x69, 0x4c, 0xc3, 0x82, 0xf1, 0xd4, 0x08, 0x02, 0x55, 0x11,
	0x97, 0x06, 0xc2, 0x0d, 0x96, 0xff, 0xde, 0x82, 0x6e, 0xe3, 0x26, 0xe8, 0x31, 0x38, 0xd7, 0x2c,
	0x8d, 0xd4, 0x51, 0x7b, 0x0d, 0xe5, 0x36, 0x38, 0xc3, 0xaf, 0x59, 0x1a, 0x61, 0x45, 0x43, 0xcf,
	0xa0, 0xbb, 0x60, 0x45, 0x43, 0xef, 0xd6, 0x86, 0x08, 0x5f, 0xb3, 0xa2, 0x59, 0x26, 0x58, 0x30,
	0xd3, 0x55, 0xff, 0x08, 0x1c, 0x99, 0x07, 0xb5, 0xc1, 0x7e, 0xfd, 0xf6, 0xaa, 0xbf, 0xe7, 0xbf,
	0x84, 0xde, 0x26, 0x1d, 0x7d, 0x0a, 0x20, 0xeb, 0x9a, 0xb3, 0x82, 0x8b, 0x75, 0x59, 0xbc, 0x86,
	0x07, 0xf5, 0xc1, 0x16, 0x74, 0x5e, 0x36, 0x5e, 0x2e, 0xfd, 0xdf, 0x2c, 0xb0, 0x27, 0x61, 0x8c,
	0x3e, 0x07, 0x47, 0xf0, 0x98, 0x96, 0xb7, 0x38, 0x6e, 0x16, 0x77, 0x88, 0x79, 0x4c, 0xb1, 0x82,
	0xd1, 0x3d, 0xd8, 0x5f, 0x08, 0xbe, 0xcc, 0xca, 0x14, 0xda, 0x40, 0x27, 0xd0, 0x61, 0x11, 0x4d,
	0x0b, 0x56, 0xac, 0x4b, 0x89, 0x54, 0xb6, 0x7f, 0x0a, 0x8e, 0x8c, 0x47, 0x00, 0x2d, 0x3c, 0x9b,
	0x4c, 0x67, 0xb8, 0xbf, 0x87, 0x0e, 0xc1, 0xc5, 0xb3, 0x6f, 0xbf, 0x9b, 0x5d, 0x5e, 0xcd, 0x70,
	0xdf, 0xf2, 0x7f, 0xb7, 0xc0, 0xad, 0x6a, 0x2d, 0x93, 0x99, 0x06, 0x95, 0x37, 0xa8, 0x6c, 0x34,
	0x86, 0x8e, 0x19, 0x1d, 0xea, 0x04, 0xbd, 0xf1, 0x83, 0x0f, 0xa7, 0xcc, 0x37, 0x24, 0xa1, 0xb8,
	0xe2, 0xa1, 0x47, 0xd0, 0xd6, 0xb3, 0xd0, 0x48, 0xe2, 0xa8, 0xfe, 0xac, 0x94, 0x1f, 0x1b, 0xdc,
	0xff, 0xeb, 0x0e, 0x74, 0x8c, 0x4c, 0xb6, 0x4a, 0xf0, 0x0c, 0xf6, 0x53, 0x4a, 0xa3, 0xbc, 0xdc,
	0x1c, 0x7d, 0x38, 0xa6, 0xb0, 0x26, 0xa0, 0xa1, 0x1a, 0x34, 0x2b, 0x16, 0xd1, 0x5c, 0x95, 0x64,
	0x3b, 0xb9, 0xe2, 0xa0, 0x53, 0x38, 0xc8, 0x48, 0xf1, 0x53, 0x30, 0x67, 0x71, 0x21, 0xd5, 0x2b,
	0xb5, 0xe8, 0xe2, 0xae, 0xf4, 0xbd, 0xd2, 0x2e, 0x59, 0x7b, 0x7e, 0x93, 0x52, 0xa1, 0xc6, 0x89,
	0x8b, 0xb5, 0x81, 0x3e, 0x01, 0x37, 0xe4, 0x49, 0xc6, 0x53, 0x9a, 0x16, 0x5e, 0x5b, 0x21, 0xb5,
	0x03, 0x9d, 0x43, 0x57, 0x5f, 0x2e, 0x88, 0xe8, 0x3c, 0xf7, 0x3a, 0xb7, 0x14, 0xae, 0x0b, 0x30,
	0xa5, 0x73, 0x0c, 0xa1, 0x59, 0xe6, 0xe8, 0x21, 0xec, 0xb3, 0x24, 0x8b, 0x73, 0xcf, 0x55, 0xf4,
	0xc3, 0x8a, 0xfe, 0x36, 0xc9, 0x62, 0xac, 0x31, 0xff, 0x39, 0xb8, 0x55, 0xf4, 0xd6, 0x5a, 0x79,
	0xd0, 0x8e, 0xe8, 0x9c, 0x2c, 0x63, 0x33, 0x68, 0x8c, 0xe9, 0xbf, 0xb7, 0xc1, 0x91, 0xa9, 0xd0,
	0x14, 0x90, 0xaa, 0x56, 0x30, 0xe7, 0x22, 0xa8, 0x1a, 0x6b, 0xed, 0x6c, 0x6c, 0x5f, 0x45, 0xbc,
	0xe2, 0xc2, 0xb8, 0xd1, 0x57, 0x70, 0xdf, 0x94, 0x71, 0x33, 0xd1, 0x6e, 0x85, 0xdc, 0x35, 0x41,
	0xcd, 0x5c, 0x13, 0xe8, 0x8b, 0x65, 0x5a, 0xb0, 0x84, 0xd6, 0x69, 0xec, 0x9d, 0x69, 0x8e, 0x4a,
	0x7e, 0x95, 0xe2, 0x39, 0x1c, 0x86, 0x2c, 0x8b, 0x82, 0xcc, 0x0c, 0x73, 0x67, 0xc7, 0x30, 0x3f,
	0x08, 0x6b, 0x23, 0x47, 0x8f, 0xa0, 0xa5, 0x87, 0xbb, 0x7a, 0x59, 0x9a, 0x4a, 0xd5, 0xef, 0xc3,
	0x9b, 0x3d, 0x5c, 0x12, 0xd0, 0x00, 0xec, 0x70, 0xfb, 0xcb, 0xf2, 0x66, 0x0f, 0x4b, 0x48, 0x7e,
	0x47, 0x11, 0x25, 0x51, 0xcc, 0x52, 0xaa, 0x74, 0xb1, 0x8f, 0x2b, 0xfb, 0x65, 0x0b, 0x1c, 0xd9,
	0x45, 0xff, 0x4f, 0x0b, 0x5a, 0x3a, 0xf5, 0x7f, 0x8e, 0x0e, 0x04, 0x8e, 0x14, 0x63, 0xd9, 0x4b,
	0xb5, 0x96, 0x5b, 0x08, 0xba, 0x62, 0x39, 0xe3, 0xa9, 0xf9, 0xee, 0x8d, 0x8d, 0x9e, 0x00, 0x64,
	0x82, 0x67, 0x54, 0x14, 0xac, 0xaa, 0xc1, 0x71, 0xf3, 0xad, 0x95, 0xd0, 0x1a, 0x37, 0x48, 0xfe,
	0x18, 0x3a, 0xc6, 0x2f, 0x27, 0xd5, 0x35, 0x35, 0xe7, 0x90, 0x4b, 0x29, 0xff, 0x15, 0x89, 0x97,
	0xd4, 0x8c, 0x1e, 0x65, 0xf8, 0x63, 0x68, 0xed, 0x78, 0x32, 0xb6, 0xc7, 0x3c, 0x06, 0x5b, 0x3e,
	0xbb, 0x08, 0x1c, 0xfa, 0x0b, 0x0d, 0x4d, 0x80, 0x5c, 0x4b, 0x1f, 0x11, 0x0b, 0xfd, 0xc6, 0xb8,
	0x58, 0xad, 0xfd, 0x1f, 0xa0, 0xdb, 0x68, 0x99, 0xfe, 0x52, 0xd5, 0x32, 0x68, 0xec, 0xd7, 0x2d,
	0x7d, 0x52, 0x13, 0x5b, 0x6b, 0xe5, 0x41, 0x7b, 0x45, 0x45, 0xa3, 0x54, 0xc6, 0xfc, 0xb1, 0xa5,
	0x7e, 0xaa, 0xce, 0xff, 0x0d, 0x00, 0x00, 0xff, 0xff, 0x10, 0x15, 0xa1, 0xeb, 0xc8, 0x09, 0x00,
	0x00,
}
