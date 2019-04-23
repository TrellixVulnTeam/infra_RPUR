// Code generated by protoc-gen-go. DO NOT EDIT.
// source: infra/qscheduler/qslib/protos/metrics/metrics.proto

package metrics

import (
	fmt "fmt"
	math "math"

	proto "github.com/golang/protobuf/proto"
	timestamp "github.com/golang/protobuf/ptypes/timestamp"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type TaskEvent_EventCategory int32

const (
	TaskEvent_CATEGORY_UNSPECIFIED TaskEvent_EventCategory = 0
	// CATEGORY_SWARMING events are due to new information from swarming.
	TaskEvent_CATEGORY_SWARMING TaskEvent_EventCategory = 256
	// CATEGORY_QSCHEDULER events are due to a decision made by quotascheduler.
	TaskEvent_CATEGORY_QSCHEDULER TaskEvent_EventCategory = 512
)

var TaskEvent_EventCategory_name = map[int32]string{
	0:   "CATEGORY_UNSPECIFIED",
	256: "CATEGORY_SWARMING",
	512: "CATEGORY_QSCHEDULER",
}

var TaskEvent_EventCategory_value = map[string]int32{
	"CATEGORY_UNSPECIFIED": 0,
	"CATEGORY_SWARMING":    256,
	"CATEGORY_QSCHEDULER":  512,
}

func (x TaskEvent_EventCategory) String() string {
	return proto.EnumName(TaskEvent_EventCategory_name, int32(x))
}

func (TaskEvent_EventCategory) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_ee70d48c7a08b27b, []int{0, 0}
}

type TaskEvent_EventType int32

const (
	// Invalid or unspecified event type.
	TaskEvent_UNSPECIFIED TaskEvent_EventType = 0
	// Task was enqueued.
	TaskEvent_SWARMING_ENQUEUED TaskEvent_EventType = 257
	// Task completed (was removed from running worker or from queue).
	TaskEvent_SWARMING_COMPLETED TaskEvent_EventType = 258
	// Task was assigned to a bot.
	TaskEvent_QSCHEDULER_ASSIGNED TaskEvent_EventType = 513
	// Task (which was previously assigned to a bot) was interrupted by another
	// task.
	TaskEvent_QSCHEDULER_PREEMPTED TaskEvent_EventType = 514
	// Task (which was previously assigned to a bot) changed its running
	// priority.
	TaskEvent_QSCHEDULER_REPRIORITIZED TaskEvent_EventType = 515
	// Task (which was previously assigned to a bot) was unassigned and moved
	// back to the queue.
	TaskEvent_QSCHEDULER_UNASSIGNED TaskEvent_EventType = 516
)

var TaskEvent_EventType_name = map[int32]string{
	0:   "UNSPECIFIED",
	257: "SWARMING_ENQUEUED",
	258: "SWARMING_COMPLETED",
	513: "QSCHEDULER_ASSIGNED",
	514: "QSCHEDULER_PREEMPTED",
	515: "QSCHEDULER_REPRIORITIZED",
	516: "QSCHEDULER_UNASSIGNED",
}

var TaskEvent_EventType_value = map[string]int32{
	"UNSPECIFIED":              0,
	"SWARMING_ENQUEUED":        257,
	"SWARMING_COMPLETED":       258,
	"QSCHEDULER_ASSIGNED":      513,
	"QSCHEDULER_PREEMPTED":     514,
	"QSCHEDULER_REPRIORITIZED": 515,
	"QSCHEDULER_UNASSIGNED":    516,
}

func (x TaskEvent_EventType) String() string {
	return proto.EnumName(TaskEvent_EventType_name, int32(x))
}

func (TaskEvent_EventType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_ee70d48c7a08b27b, []int{0, 1}
}

type TaskEvent_CompletedDetails_Reason int32

const (
	TaskEvent_CompletedDetails_UNSPECIFIED TaskEvent_CompletedDetails_Reason = 0
	// The bot that was previously assigned this task became idle.
	TaskEvent_CompletedDetails_BOT_IDLE TaskEvent_CompletedDetails_Reason = 1
	// The task became absent while it was still idle.
	TaskEvent_CompletedDetails_IDLE_TASK_ABSENT TaskEvent_CompletedDetails_Reason = 2
	// The task becamse absent while it was running.
	TaskEvent_CompletedDetails_RUNNING_TASK_ABSENT TaskEvent_CompletedDetails_Reason = 3
	// The task was believed to be running on a bot, but the scheduler was
	// notified that it was running on a different bot, so its state was
	// deleted.
	TaskEvent_CompletedDetails_INCONSISTENT_BOT_FOR_TASK TaskEvent_CompletedDetails_Reason = 4
	// The task was believed the be running on a bot, but the scheduler was
	// notified that a different task is running on the bot, so its state was
	// deleted.
	TaskEvent_CompletedDetails_INCONSISTENT_TASK_FOR_BOT TaskEvent_CompletedDetails_Reason = 5
)

var TaskEvent_CompletedDetails_Reason_name = map[int32]string{
	0: "UNSPECIFIED",
	1: "BOT_IDLE",
	2: "IDLE_TASK_ABSENT",
	3: "RUNNING_TASK_ABSENT",
	4: "INCONSISTENT_BOT_FOR_TASK",
	5: "INCONSISTENT_TASK_FOR_BOT",
}

var TaskEvent_CompletedDetails_Reason_value = map[string]int32{
	"UNSPECIFIED":               0,
	"BOT_IDLE":                  1,
	"IDLE_TASK_ABSENT":          2,
	"RUNNING_TASK_ABSENT":       3,
	"INCONSISTENT_BOT_FOR_TASK": 4,
	"INCONSISTENT_TASK_FOR_BOT": 5,
}

func (x TaskEvent_CompletedDetails_Reason) String() string {
	return proto.EnumName(TaskEvent_CompletedDetails_Reason_name, int32(x))
}

func (TaskEvent_CompletedDetails_Reason) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_ee70d48c7a08b27b, []int{0, 4, 0}
}

// TaskEvent represents a quotascheduler event that happened to a particular
// task, for metrics purposes.
//
// This proto is intended to be used as a schema for a BigQuery table, in which
// events are logged.
type TaskEvent struct {
	// EventType is the type of event that occurred.
	EventType TaskEvent_EventType `protobuf:"varint,1,opt,name=event_type,json=eventType,proto3,enum=metrics.TaskEvent_EventType" json:"event_type,omitempty"`
	// SchedulerId is the ID of the scheduler in which the event occurred.
	SchedulerId string `protobuf:"bytes,2,opt,name=scheduler_id,json=schedulerId,proto3" json:"scheduler_id,omitempty"`
	// TaskId is the task ID that the event happened to.
	TaskId string `protobuf:"bytes,3,opt,name=task_id,json=taskId,proto3" json:"task_id,omitempty"`
	// Time is the time at which the event happened.
	Time *timestamp.Timestamp `protobuf:"bytes,4,opt,name=time,proto3" json:"time,omitempty"`
	// BaseLabels are the base labels of the task.
	BaseLabels []string `protobuf:"bytes,5,rep,name=base_labels,json=baseLabels,proto3" json:"base_labels,omitempty"`
	// ProvisionableLabels are the provisionable labels of the task.
	ProvisionableLabels []string `protobuf:"bytes,6,rep,name=provisionable_labels,json=provisionableLabels,proto3" json:"provisionable_labels,omitempty"`
	// AccountId is the quotascheduler account that the task will be charged to.
	AccountId string `protobuf:"bytes,7,opt,name=account_id,json=accountId,proto3" json:"account_id,omitempty"`
	// AccountValid indicates whether this task's quotascheduler account is valid
	// (exists) at this time.
	AccountValid bool `protobuf:"varint,8,opt,name=account_valid,json=accountValid,proto3" json:"account_valid,omitempty"`
	// AccountBalance is the task's quotascheduler account's balance at this time.
	AccountBalance []float32 `protobuf:"fixed32,9,rep,packed,name=account_balance,json=accountBalance,proto3" json:"account_balance,omitempty"`
	// Cost is the total quota cost accumulated so far to the quotascheduler
	// account, when running this task.
	Cost []float32 `protobuf:"fixed32,10,rep,packed,name=cost,proto3" json:"cost,omitempty"`
	// BotId is the bot that the event occurred on (relevant for all event
	// types except for ENQUEUED).
	BotId string `protobuf:"bytes,11,opt,name=bot_id,json=botId,proto3" json:"bot_id,omitempty"`
	// BotDimensions are the dimensions of the bot (if relevant).
	BotDimensions []string `protobuf:"bytes,12,rep,name=bot_dimensions,json=botDimensions,proto3" json:"bot_dimensions,omitempty"`
	// Category is the EventCategory for this event's EventType.
	Category TaskEvent_EventCategory `protobuf:"varint,13,opt,name=category,proto3,enum=metrics.TaskEvent_EventCategory" json:"category,omitempty"`
	// IsCallback specifies whether this event was due to a notification callback
	// from swarming that was requested by the scheduler. This is for diagnostic
	// purposes only, and is only applicable for events of category
	// CATEGORY_SWARMING.
	IsCallback bool `protobuf:"varint,14,opt,name=is_callback,json=isCallback,proto3" json:"is_callback,omitempty"`
	// Types that are valid to be assigned to Details:
	//	*TaskEvent_EnqueuedDetails_
	//	*TaskEvent_AssignedDetails_
	//	*TaskEvent_PreemptedDetails_
	//	*TaskEvent_ReprioritizedDetails_
	//	*TaskEvent_CompletedDetails_
	//	*TaskEvent_UnassignedDetails_
	Details              isTaskEvent_Details `protobuf_oneof:"details"`
	XXX_NoUnkeyedLiteral struct{}            `json:"-"`
	XXX_unrecognized     []byte              `json:"-"`
	XXX_sizecache        int32               `json:"-"`
}

func (m *TaskEvent) Reset()         { *m = TaskEvent{} }
func (m *TaskEvent) String() string { return proto.CompactTextString(m) }
func (*TaskEvent) ProtoMessage()    {}
func (*TaskEvent) Descriptor() ([]byte, []int) {
	return fileDescriptor_ee70d48c7a08b27b, []int{0}
}

func (m *TaskEvent) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_TaskEvent.Unmarshal(m, b)
}
func (m *TaskEvent) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_TaskEvent.Marshal(b, m, deterministic)
}
func (m *TaskEvent) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TaskEvent.Merge(m, src)
}
func (m *TaskEvent) XXX_Size() int {
	return xxx_messageInfo_TaskEvent.Size(m)
}
func (m *TaskEvent) XXX_DiscardUnknown() {
	xxx_messageInfo_TaskEvent.DiscardUnknown(m)
}

var xxx_messageInfo_TaskEvent proto.InternalMessageInfo

func (m *TaskEvent) GetEventType() TaskEvent_EventType {
	if m != nil {
		return m.EventType
	}
	return TaskEvent_UNSPECIFIED
}

func (m *TaskEvent) GetSchedulerId() string {
	if m != nil {
		return m.SchedulerId
	}
	return ""
}

func (m *TaskEvent) GetTaskId() string {
	if m != nil {
		return m.TaskId
	}
	return ""
}

func (m *TaskEvent) GetTime() *timestamp.Timestamp {
	if m != nil {
		return m.Time
	}
	return nil
}

func (m *TaskEvent) GetBaseLabels() []string {
	if m != nil {
		return m.BaseLabels
	}
	return nil
}

func (m *TaskEvent) GetProvisionableLabels() []string {
	if m != nil {
		return m.ProvisionableLabels
	}
	return nil
}

func (m *TaskEvent) GetAccountId() string {
	if m != nil {
		return m.AccountId
	}
	return ""
}

func (m *TaskEvent) GetAccountValid() bool {
	if m != nil {
		return m.AccountValid
	}
	return false
}

func (m *TaskEvent) GetAccountBalance() []float32 {
	if m != nil {
		return m.AccountBalance
	}
	return nil
}

func (m *TaskEvent) GetCost() []float32 {
	if m != nil {
		return m.Cost
	}
	return nil
}

func (m *TaskEvent) GetBotId() string {
	if m != nil {
		return m.BotId
	}
	return ""
}

func (m *TaskEvent) GetBotDimensions() []string {
	if m != nil {
		return m.BotDimensions
	}
	return nil
}

func (m *TaskEvent) GetCategory() TaskEvent_EventCategory {
	if m != nil {
		return m.Category
	}
	return TaskEvent_CATEGORY_UNSPECIFIED
}

func (m *TaskEvent) GetIsCallback() bool {
	if m != nil {
		return m.IsCallback
	}
	return false
}

type isTaskEvent_Details interface {
	isTaskEvent_Details()
}

type TaskEvent_EnqueuedDetails_ struct {
	EnqueuedDetails *TaskEvent_EnqueuedDetails `protobuf:"bytes,100,opt,name=enqueued_details,json=enqueuedDetails,proto3,oneof"`
}

type TaskEvent_AssignedDetails_ struct {
	AssignedDetails *TaskEvent_AssignedDetails `protobuf:"bytes,101,opt,name=assigned_details,json=assignedDetails,proto3,oneof"`
}

type TaskEvent_PreemptedDetails_ struct {
	PreemptedDetails *TaskEvent_PreemptedDetails `protobuf:"bytes,102,opt,name=preempted_details,json=preemptedDetails,proto3,oneof"`
}

type TaskEvent_ReprioritizedDetails_ struct {
	ReprioritizedDetails *TaskEvent_ReprioritizedDetails `protobuf:"bytes,103,opt,name=reprioritized_details,json=reprioritizedDetails,proto3,oneof"`
}

type TaskEvent_CompletedDetails_ struct {
	CompletedDetails *TaskEvent_CompletedDetails `protobuf:"bytes,104,opt,name=completed_details,json=completedDetails,proto3,oneof"`
}

type TaskEvent_UnassignedDetails_ struct {
	UnassignedDetails *TaskEvent_UnassignedDetails `protobuf:"bytes,105,opt,name=unassigned_details,json=unassignedDetails,proto3,oneof"`
}

func (*TaskEvent_EnqueuedDetails_) isTaskEvent_Details() {}

func (*TaskEvent_AssignedDetails_) isTaskEvent_Details() {}

func (*TaskEvent_PreemptedDetails_) isTaskEvent_Details() {}

func (*TaskEvent_ReprioritizedDetails_) isTaskEvent_Details() {}

func (*TaskEvent_CompletedDetails_) isTaskEvent_Details() {}

func (*TaskEvent_UnassignedDetails_) isTaskEvent_Details() {}

func (m *TaskEvent) GetDetails() isTaskEvent_Details {
	if m != nil {
		return m.Details
	}
	return nil
}

func (m *TaskEvent) GetEnqueuedDetails() *TaskEvent_EnqueuedDetails {
	if x, ok := m.GetDetails().(*TaskEvent_EnqueuedDetails_); ok {
		return x.EnqueuedDetails
	}
	return nil
}

func (m *TaskEvent) GetAssignedDetails() *TaskEvent_AssignedDetails {
	if x, ok := m.GetDetails().(*TaskEvent_AssignedDetails_); ok {
		return x.AssignedDetails
	}
	return nil
}

func (m *TaskEvent) GetPreemptedDetails() *TaskEvent_PreemptedDetails {
	if x, ok := m.GetDetails().(*TaskEvent_PreemptedDetails_); ok {
		return x.PreemptedDetails
	}
	return nil
}

func (m *TaskEvent) GetReprioritizedDetails() *TaskEvent_ReprioritizedDetails {
	if x, ok := m.GetDetails().(*TaskEvent_ReprioritizedDetails_); ok {
		return x.ReprioritizedDetails
	}
	return nil
}

func (m *TaskEvent) GetCompletedDetails() *TaskEvent_CompletedDetails {
	if x, ok := m.GetDetails().(*TaskEvent_CompletedDetails_); ok {
		return x.CompletedDetails
	}
	return nil
}

func (m *TaskEvent) GetUnassignedDetails() *TaskEvent_UnassignedDetails {
	if x, ok := m.GetDetails().(*TaskEvent_UnassignedDetails_); ok {
		return x.UnassignedDetails
	}
	return nil
}

// XXX_OneofWrappers is for the internal use of the proto package.
func (*TaskEvent) XXX_OneofWrappers() []interface{} {
	return []interface{}{
		(*TaskEvent_EnqueuedDetails_)(nil),
		(*TaskEvent_AssignedDetails_)(nil),
		(*TaskEvent_PreemptedDetails_)(nil),
		(*TaskEvent_ReprioritizedDetails_)(nil),
		(*TaskEvent_CompletedDetails_)(nil),
		(*TaskEvent_UnassignedDetails_)(nil),
	}
}

// EnqueuedDetails represents event details that are specific to the
// ENQUEUED event type.
type TaskEvent_EnqueuedDetails struct {
	// Tags is the set of swarming tags for the task.
	Tags                 []string `protobuf:"bytes,1,rep,name=tags,proto3" json:"tags,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *TaskEvent_EnqueuedDetails) Reset()         { *m = TaskEvent_EnqueuedDetails{} }
func (m *TaskEvent_EnqueuedDetails) String() string { return proto.CompactTextString(m) }
func (*TaskEvent_EnqueuedDetails) ProtoMessage()    {}
func (*TaskEvent_EnqueuedDetails) Descriptor() ([]byte, []int) {
	return fileDescriptor_ee70d48c7a08b27b, []int{0, 0}
}

func (m *TaskEvent_EnqueuedDetails) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_TaskEvent_EnqueuedDetails.Unmarshal(m, b)
}
func (m *TaskEvent_EnqueuedDetails) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_TaskEvent_EnqueuedDetails.Marshal(b, m, deterministic)
}
func (m *TaskEvent_EnqueuedDetails) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TaskEvent_EnqueuedDetails.Merge(m, src)
}
func (m *TaskEvent_EnqueuedDetails) XXX_Size() int {
	return xxx_messageInfo_TaskEvent_EnqueuedDetails.Size(m)
}
func (m *TaskEvent_EnqueuedDetails) XXX_DiscardUnknown() {
	xxx_messageInfo_TaskEvent_EnqueuedDetails.DiscardUnknown(m)
}

var xxx_messageInfo_TaskEvent_EnqueuedDetails proto.InternalMessageInfo

func (m *TaskEvent_EnqueuedDetails) GetTags() []string {
	if m != nil {
		return m.Tags
	}
	return nil
}

// AssignedDetails represents event details that are specific to the
// ASSIGNED event type.
type TaskEvent_AssignedDetails struct {
	// ProvisionRequired is whether provision is required to run this task
	// on the bot (i.e. if a slice other than the 0th slice was selected).
	ProvisionRequired bool `protobuf:"varint,1,opt,name=provision_required,json=provisionRequired,proto3" json:"provision_required,omitempty"`
	// Priority is the qscheduler priority that the task is running at.
	Priority int32 `protobuf:"varint,2,opt,name=priority,proto3" json:"priority,omitempty"`
	// Preempting is true if this task preempted another one that was already
	// running on the bot.
	Preempting bool `protobuf:"varint,3,opt,name=preempting,proto3" json:"preempting,omitempty"`
	// PreemptionCost is the cost paid by this task's account in order to
	// preempt the previous task on this bot, if this was a preempting
	// assignment.
	PreemptionCost []float32 `protobuf:"fixed32,4,rep,packed,name=preemption_cost,json=preemptionCost,proto3" json:"preemption_cost,omitempty"`
	// PreemptedTaskId is the task that was preempted, if any.
	PreemptedTaskId      string   `protobuf:"bytes,5,opt,name=preempted_task_id,json=preemptedTaskId,proto3" json:"preempted_task_id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *TaskEvent_AssignedDetails) Reset()         { *m = TaskEvent_AssignedDetails{} }
func (m *TaskEvent_AssignedDetails) String() string { return proto.CompactTextString(m) }
func (*TaskEvent_AssignedDetails) ProtoMessage()    {}
func (*TaskEvent_AssignedDetails) Descriptor() ([]byte, []int) {
	return fileDescriptor_ee70d48c7a08b27b, []int{0, 1}
}

func (m *TaskEvent_AssignedDetails) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_TaskEvent_AssignedDetails.Unmarshal(m, b)
}
func (m *TaskEvent_AssignedDetails) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_TaskEvent_AssignedDetails.Marshal(b, m, deterministic)
}
func (m *TaskEvent_AssignedDetails) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TaskEvent_AssignedDetails.Merge(m, src)
}
func (m *TaskEvent_AssignedDetails) XXX_Size() int {
	return xxx_messageInfo_TaskEvent_AssignedDetails.Size(m)
}
func (m *TaskEvent_AssignedDetails) XXX_DiscardUnknown() {
	xxx_messageInfo_TaskEvent_AssignedDetails.DiscardUnknown(m)
}

var xxx_messageInfo_TaskEvent_AssignedDetails proto.InternalMessageInfo

func (m *TaskEvent_AssignedDetails) GetProvisionRequired() bool {
	if m != nil {
		return m.ProvisionRequired
	}
	return false
}

func (m *TaskEvent_AssignedDetails) GetPriority() int32 {
	if m != nil {
		return m.Priority
	}
	return 0
}

func (m *TaskEvent_AssignedDetails) GetPreempting() bool {
	if m != nil {
		return m.Preempting
	}
	return false
}

func (m *TaskEvent_AssignedDetails) GetPreemptionCost() []float32 {
	if m != nil {
		return m.PreemptionCost
	}
	return nil
}

func (m *TaskEvent_AssignedDetails) GetPreemptedTaskId() string {
	if m != nil {
		return m.PreemptedTaskId
	}
	return ""
}

// PreemptedDetails represents event details that are specific to the
// PREEMPTED event type.
type TaskEvent_PreemptedDetails struct {
	// PreemptingAccountId is the account id of the task that preempted this
	// task.
	PreemptingAccountId string `protobuf:"bytes,1,opt,name=preempting_account_id,json=preemptingAccountId,proto3" json:"preempting_account_id,omitempty"`
	// PreemptingTaskId is the task id of the task that preempted this task.
	PreemptingTaskId string `protobuf:"bytes,2,opt,name=preempting_task_id,json=preemptingTaskId,proto3" json:"preempting_task_id,omitempty"`
	// Priority is the priority that this task was running at prior to being
	// preempted.
	Priority int32 `protobuf:"varint,3,opt,name=priority,proto3" json:"priority,omitempty"`
	// PreemptingPriority is the priority of the task that preempted this task.
	PreemptingPriority   int32    `protobuf:"varint,4,opt,name=preempting_priority,json=preemptingPriority,proto3" json:"preempting_priority,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *TaskEvent_PreemptedDetails) Reset()         { *m = TaskEvent_PreemptedDetails{} }
func (m *TaskEvent_PreemptedDetails) String() string { return proto.CompactTextString(m) }
func (*TaskEvent_PreemptedDetails) ProtoMessage()    {}
func (*TaskEvent_PreemptedDetails) Descriptor() ([]byte, []int) {
	return fileDescriptor_ee70d48c7a08b27b, []int{0, 2}
}

func (m *TaskEvent_PreemptedDetails) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_TaskEvent_PreemptedDetails.Unmarshal(m, b)
}
func (m *TaskEvent_PreemptedDetails) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_TaskEvent_PreemptedDetails.Marshal(b, m, deterministic)
}
func (m *TaskEvent_PreemptedDetails) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TaskEvent_PreemptedDetails.Merge(m, src)
}
func (m *TaskEvent_PreemptedDetails) XXX_Size() int {
	return xxx_messageInfo_TaskEvent_PreemptedDetails.Size(m)
}
func (m *TaskEvent_PreemptedDetails) XXX_DiscardUnknown() {
	xxx_messageInfo_TaskEvent_PreemptedDetails.DiscardUnknown(m)
}

var xxx_messageInfo_TaskEvent_PreemptedDetails proto.InternalMessageInfo

func (m *TaskEvent_PreemptedDetails) GetPreemptingAccountId() string {
	if m != nil {
		return m.PreemptingAccountId
	}
	return ""
}

func (m *TaskEvent_PreemptedDetails) GetPreemptingTaskId() string {
	if m != nil {
		return m.PreemptingTaskId
	}
	return ""
}

func (m *TaskEvent_PreemptedDetails) GetPriority() int32 {
	if m != nil {
		return m.Priority
	}
	return 0
}

func (m *TaskEvent_PreemptedDetails) GetPreemptingPriority() int32 {
	if m != nil {
		return m.PreemptingPriority
	}
	return 0
}

// ReprioritizedDetails represents event details that are specific to the
// PREPRIORITIZED event type.
type TaskEvent_ReprioritizedDetails struct {
	// OldPriority is the previous priority the task was running at.
	OldPriority int32 `protobuf:"varint,1,opt,name=old_priority,json=oldPriority,proto3" json:"old_priority,omitempty"`
	// NewPriority is the new priority the task is running at.
	NewPriority          int32    `protobuf:"varint,2,opt,name=new_priority,json=newPriority,proto3" json:"new_priority,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *TaskEvent_ReprioritizedDetails) Reset()         { *m = TaskEvent_ReprioritizedDetails{} }
func (m *TaskEvent_ReprioritizedDetails) String() string { return proto.CompactTextString(m) }
func (*TaskEvent_ReprioritizedDetails) ProtoMessage()    {}
func (*TaskEvent_ReprioritizedDetails) Descriptor() ([]byte, []int) {
	return fileDescriptor_ee70d48c7a08b27b, []int{0, 3}
}

func (m *TaskEvent_ReprioritizedDetails) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_TaskEvent_ReprioritizedDetails.Unmarshal(m, b)
}
func (m *TaskEvent_ReprioritizedDetails) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_TaskEvent_ReprioritizedDetails.Marshal(b, m, deterministic)
}
func (m *TaskEvent_ReprioritizedDetails) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TaskEvent_ReprioritizedDetails.Merge(m, src)
}
func (m *TaskEvent_ReprioritizedDetails) XXX_Size() int {
	return xxx_messageInfo_TaskEvent_ReprioritizedDetails.Size(m)
}
func (m *TaskEvent_ReprioritizedDetails) XXX_DiscardUnknown() {
	xxx_messageInfo_TaskEvent_ReprioritizedDetails.DiscardUnknown(m)
}

var xxx_messageInfo_TaskEvent_ReprioritizedDetails proto.InternalMessageInfo

func (m *TaskEvent_ReprioritizedDetails) GetOldPriority() int32 {
	if m != nil {
		return m.OldPriority
	}
	return 0
}

func (m *TaskEvent_ReprioritizedDetails) GetNewPriority() int32 {
	if m != nil {
		return m.NewPriority
	}
	return 0
}

// ReprioritizedDetails represents event details that are specific to the
// COMPLETED event type.
type TaskEvent_CompletedDetails struct {
	Reason TaskEvent_CompletedDetails_Reason `protobuf:"varint,1,opt,name=reason,proto3,enum=metrics.TaskEvent_CompletedDetails_Reason" json:"reason,omitempty"`
	// If reason is INCONSISTENT_TASK_FOR_BOT, other_task specifies which task
	// it was that notified itself on this bot.
	OtherTask string `protobuf:"bytes,2,opt,name=other_task,json=otherTask,proto3" json:"other_task,omitempty"`
	// If reason is INCONSISTENT_BOT_FOR_TASK, other_bot specifies which bot
	// it was that was notified with this task.
	OtherBot             string   `protobuf:"bytes,3,opt,name=other_bot,json=otherBot,proto3" json:"other_bot,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *TaskEvent_CompletedDetails) Reset()         { *m = TaskEvent_CompletedDetails{} }
func (m *TaskEvent_CompletedDetails) String() string { return proto.CompactTextString(m) }
func (*TaskEvent_CompletedDetails) ProtoMessage()    {}
func (*TaskEvent_CompletedDetails) Descriptor() ([]byte, []int) {
	return fileDescriptor_ee70d48c7a08b27b, []int{0, 4}
}

func (m *TaskEvent_CompletedDetails) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_TaskEvent_CompletedDetails.Unmarshal(m, b)
}
func (m *TaskEvent_CompletedDetails) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_TaskEvent_CompletedDetails.Marshal(b, m, deterministic)
}
func (m *TaskEvent_CompletedDetails) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TaskEvent_CompletedDetails.Merge(m, src)
}
func (m *TaskEvent_CompletedDetails) XXX_Size() int {
	return xxx_messageInfo_TaskEvent_CompletedDetails.Size(m)
}
func (m *TaskEvent_CompletedDetails) XXX_DiscardUnknown() {
	xxx_messageInfo_TaskEvent_CompletedDetails.DiscardUnknown(m)
}

var xxx_messageInfo_TaskEvent_CompletedDetails proto.InternalMessageInfo

func (m *TaskEvent_CompletedDetails) GetReason() TaskEvent_CompletedDetails_Reason {
	if m != nil {
		return m.Reason
	}
	return TaskEvent_CompletedDetails_UNSPECIFIED
}

func (m *TaskEvent_CompletedDetails) GetOtherTask() string {
	if m != nil {
		return m.OtherTask
	}
	return ""
}

func (m *TaskEvent_CompletedDetails) GetOtherBot() string {
	if m != nil {
		return m.OtherBot
	}
	return ""
}

// UnassignedDetails represents event details that are specific to the
// UNASSIGNED event type.
type TaskEvent_UnassignedDetails struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *TaskEvent_UnassignedDetails) Reset()         { *m = TaskEvent_UnassignedDetails{} }
func (m *TaskEvent_UnassignedDetails) String() string { return proto.CompactTextString(m) }
func (*TaskEvent_UnassignedDetails) ProtoMessage()    {}
func (*TaskEvent_UnassignedDetails) Descriptor() ([]byte, []int) {
	return fileDescriptor_ee70d48c7a08b27b, []int{0, 5}
}

func (m *TaskEvent_UnassignedDetails) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_TaskEvent_UnassignedDetails.Unmarshal(m, b)
}
func (m *TaskEvent_UnassignedDetails) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_TaskEvent_UnassignedDetails.Marshal(b, m, deterministic)
}
func (m *TaskEvent_UnassignedDetails) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TaskEvent_UnassignedDetails.Merge(m, src)
}
func (m *TaskEvent_UnassignedDetails) XXX_Size() int {
	return xxx_messageInfo_TaskEvent_UnassignedDetails.Size(m)
}
func (m *TaskEvent_UnassignedDetails) XXX_DiscardUnknown() {
	xxx_messageInfo_TaskEvent_UnassignedDetails.DiscardUnknown(m)
}

var xxx_messageInfo_TaskEvent_UnassignedDetails proto.InternalMessageInfo

// EventList wraps a repeated TaskEvent.
type EventList struct {
	Events               []*TaskEvent `protobuf:"bytes,1,rep,name=events,proto3" json:"events,omitempty"`
	XXX_NoUnkeyedLiteral struct{}     `json:"-"`
	XXX_unrecognized     []byte       `json:"-"`
	XXX_sizecache        int32        `json:"-"`
}

func (m *EventList) Reset()         { *m = EventList{} }
func (m *EventList) String() string { return proto.CompactTextString(m) }
func (*EventList) ProtoMessage()    {}
func (*EventList) Descriptor() ([]byte, []int) {
	return fileDescriptor_ee70d48c7a08b27b, []int{1}
}

func (m *EventList) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_EventList.Unmarshal(m, b)
}
func (m *EventList) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_EventList.Marshal(b, m, deterministic)
}
func (m *EventList) XXX_Merge(src proto.Message) {
	xxx_messageInfo_EventList.Merge(m, src)
}
func (m *EventList) XXX_Size() int {
	return xxx_messageInfo_EventList.Size(m)
}
func (m *EventList) XXX_DiscardUnknown() {
	xxx_messageInfo_EventList.DiscardUnknown(m)
}

var xxx_messageInfo_EventList proto.InternalMessageInfo

func (m *EventList) GetEvents() []*TaskEvent {
	if m != nil {
		return m.Events
	}
	return nil
}

func init() {
	proto.RegisterEnum("metrics.TaskEvent_EventCategory", TaskEvent_EventCategory_name, TaskEvent_EventCategory_value)
	proto.RegisterEnum("metrics.TaskEvent_EventType", TaskEvent_EventType_name, TaskEvent_EventType_value)
	proto.RegisterEnum("metrics.TaskEvent_CompletedDetails_Reason", TaskEvent_CompletedDetails_Reason_name, TaskEvent_CompletedDetails_Reason_value)
	proto.RegisterType((*TaskEvent)(nil), "metrics.TaskEvent")
	proto.RegisterType((*TaskEvent_EnqueuedDetails)(nil), "metrics.TaskEvent.EnqueuedDetails")
	proto.RegisterType((*TaskEvent_AssignedDetails)(nil), "metrics.TaskEvent.AssignedDetails")
	proto.RegisterType((*TaskEvent_PreemptedDetails)(nil), "metrics.TaskEvent.PreemptedDetails")
	proto.RegisterType((*TaskEvent_ReprioritizedDetails)(nil), "metrics.TaskEvent.ReprioritizedDetails")
	proto.RegisterType((*TaskEvent_CompletedDetails)(nil), "metrics.TaskEvent.CompletedDetails")
	proto.RegisterType((*TaskEvent_UnassignedDetails)(nil), "metrics.TaskEvent.UnassignedDetails")
	proto.RegisterType((*EventList)(nil), "metrics.EventList")
}

func init() {
	proto.RegisterFile("infra/qscheduler/qslib/protos/metrics/metrics.proto", fileDescriptor_ee70d48c7a08b27b)
}

var fileDescriptor_ee70d48c7a08b27b = []byte{
	// 1093 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x56, 0xcb, 0x6e, 0xdb, 0x46,
	0x17, 0x0e, 0x6d, 0x5a, 0x96, 0x8e, 0x6c, 0x8b, 0x1a, 0xcb, 0x09, 0xa3, 0xff, 0x4f, 0xa3, 0x2a,
	0x0d, 0x22, 0x04, 0xad, 0x84, 0x3a, 0x8b, 0x2e, 0xda, 0x8d, 0x2e, 0x8c, 0xc3, 0xd6, 0xa6, 0x94,
	0x11, 0xd5, 0xa2, 0x17, 0x94, 0xe0, 0x65, 0x2c, 0x13, 0xa6, 0x38, 0x32, 0x49, 0x39, 0x70, 0x57,
	0xe9, 0xe5, 0x11, 0xfa, 0x12, 0x5d, 0x77, 0xdd, 0x45, 0x1f, 0xa1, 0x6f, 0x54, 0xcc, 0xf0, 0x22,
	0xea, 0x62, 0xa0, 0x1b, 0x9b, 0xf3, 0x7d, 0xdf, 0xf9, 0xce, 0xe1, 0x99, 0x99, 0x43, 0xc1, 0x2b,
	0xd7, 0xbf, 0x0c, 0xcc, 0xce, 0x4d, 0x68, 0x5f, 0x11, 0x67, 0xe1, 0x91, 0xa0, 0x73, 0x13, 0x7a,
	0xae, 0xd5, 0x99, 0x07, 0x34, 0xa2, 0x61, 0x67, 0x46, 0xa2, 0xc0, 0xb5, 0xb3, 0xff, 0x6d, 0x0e,
	0xa3, 0xfd, 0x64, 0x59, 0x7f, 0x3a, 0xa5, 0x74, 0xea, 0x91, 0x58, 0x6d, 0x2d, 0x2e, 0x3b, 0x91,
	0x3b, 0x23, 0x61, 0x64, 0xce, 0xe6, 0xb1, 0xb2, 0xf9, 0x67, 0x15, 0x4a, 0xba, 0x19, 0x5e, 0x2b,
	0xb7, 0xc4, 0x8f, 0xd0, 0xe7, 0x00, 0x84, 0x3d, 0x18, 0xd1, 0xdd, 0x9c, 0xc8, 0x42, 0x43, 0x68,
	0x1d, 0x9d, 0xfe, 0xbf, 0x9d, 0x7a, 0x67, 0xba, 0x36, 0xff, 0xab, 0xdf, 0xcd, 0x09, 0x2e, 0x91,
	0xf4, 0x11, 0x7d, 0x08, 0x07, 0x59, 0x91, 0x86, 0xeb, 0xc8, 0x3b, 0x0d, 0xa1, 0x55, 0xc2, 0xe5,
	0x0c, 0x53, 0x1d, 0xf4, 0x08, 0xf6, 0x23, 0x33, 0xbc, 0x66, 0xec, 0x2e, 0x67, 0x0b, 0x6c, 0xa9,
	0x3a, 0xa8, 0x0d, 0x22, 0xab, 0x4c, 0x16, 0x1b, 0x42, 0xab, 0x7c, 0x5a, 0x6f, 0xc7, 0x65, 0xb7,
	0xd3, 0xb2, 0xdb, 0x7a, 0x5a, 0x36, 0xe6, 0x3a, 0xf4, 0x14, 0xca, 0x96, 0x19, 0x12, 0xc3, 0x33,
	0x2d, 0xe2, 0x85, 0xf2, 0x5e, 0x63, 0xb7, 0x55, 0xc2, 0xc0, 0xa0, 0x73, 0x8e, 0xa0, 0x4f, 0xa1,
	0x36, 0x0f, 0xe8, 0xad, 0x1b, 0xba, 0xd4, 0x37, 0x2d, 0x2f, 0x53, 0x16, 0xb8, 0xf2, 0x78, 0x85,
	0x4b, 0x42, 0x9e, 0x00, 0x98, 0xb6, 0x4d, 0x17, 0x7e, 0xc4, 0xea, 0xdb, 0xe7, 0xf5, 0x95, 0x12,
	0x44, 0x75, 0xd0, 0x33, 0x38, 0x4c, 0xe9, 0x5b, 0xd3, 0x73, 0x1d, 0xb9, 0xd8, 0x10, 0x5a, 0x45,
	0x7c, 0x90, 0x80, 0x5f, 0x33, 0x0c, 0xbd, 0x80, 0x4a, 0x2a, 0xb2, 0x4c, 0xcf, 0xf4, 0x6d, 0x22,
	0x97, 0x1a, 0xbb, 0xad, 0x1d, 0x7c, 0x94, 0xc0, 0xbd, 0x18, 0x45, 0x08, 0x44, 0x9b, 0x86, 0x91,
	0x0c, 0x9c, 0xe5, 0xcf, 0xe8, 0x04, 0x0a, 0x16, 0xe5, 0xc9, 0xcb, 0x3c, 0xf9, 0x9e, 0x45, 0x59,
	0xe2, 0xe7, 0x70, 0xc4, 0x60, 0xc7, 0x9d, 0x11, 0x9f, 0x95, 0x1c, 0xca, 0x07, 0xfc, 0x25, 0x0e,
	0x2d, 0x1a, 0x0d, 0x32, 0x10, 0x7d, 0x01, 0x45, 0xdb, 0x8c, 0xc8, 0x94, 0x06, 0x77, 0xf2, 0x21,
	0xdf, 0xb9, 0xc6, 0x7d, 0x3b, 0xd7, 0x4f, 0x74, 0x38, 0x8b, 0x60, 0x0d, 0x75, 0x43, 0xc3, 0x36,
	0x3d, 0xcf, 0x32, 0xed, 0x6b, 0xf9, 0x88, 0xbf, 0x1b, 0xb8, 0x61, 0x3f, 0x41, 0xd0, 0x10, 0x24,
	0xe2, 0xdf, 0x2c, 0xc8, 0x82, 0x38, 0x86, 0x43, 0x22, 0xd3, 0xf5, 0x42, 0xd9, 0xe1, 0xbb, 0xd5,
	0xdc, 0x96, 0x26, 0x91, 0x0e, 0x62, 0xe5, 0x9b, 0x07, 0xb8, 0x42, 0x56, 0x21, 0x66, 0x68, 0x86,
	0xa1, 0x3b, 0xf5, 0x73, 0x86, 0xe4, 0x5e, 0xc3, 0x6e, 0x22, 0xcd, 0x19, 0x9a, 0xab, 0x10, 0xc2,
	0x50, 0x9d, 0x07, 0x84, 0xcc, 0xe6, 0x51, 0xce, 0xf1, 0x92, 0x3b, 0x3e, 0xdb, 0xe2, 0x38, 0x4a,
	0xb5, 0x4b, 0x4b, 0x69, 0xbe, 0x86, 0xa1, 0x1f, 0xe1, 0x24, 0x20, 0xf3, 0xc0, 0xa5, 0x81, 0x1b,
	0xb9, 0x3f, 0xe5, 0x7c, 0xa7, 0xdc, 0xf7, 0xc5, 0x16, 0x5f, 0x9c, 0xd7, 0x2f, 0xbd, 0x6b, 0xc1,
	0x16, 0x9c, 0xd5, 0x6c, 0xd3, 0xd9, 0xdc, 0x23, 0xf9, 0x9a, 0xaf, 0xee, 0xad, 0xb9, 0x9f, 0x6a,
	0x73, 0x35, 0xdb, 0x6b, 0x18, 0x9a, 0x00, 0x5a, 0xf8, 0x1b, 0xad, 0x75, 0xb9, 0xe9, 0x47, 0x5b,
	0x4c, 0x27, 0xbe, 0xb9, 0xd1, 0xdc, 0xea, 0x62, 0x1d, 0xac, 0x3f, 0x87, 0xca, 0xda, 0xae, 0xb2,
	0x43, 0x1c, 0x99, 0xd3, 0x50, 0x16, 0xf8, 0x79, 0xe4, 0xcf, 0xf5, 0x7f, 0x04, 0xa8, 0xac, 0x6d,
	0x16, 0xfa, 0x04, 0x50, 0x76, 0xe1, 0x8c, 0x80, 0xdc, 0x2c, 0xdc, 0x80, 0x38, 0x7c, 0xbc, 0x14,
	0x71, 0x35, 0x63, 0x70, 0x42, 0xa0, 0x3a, 0x14, 0x93, 0x56, 0xdd, 0xf1, 0x21, 0xb2, 0x87, 0xb3,
	0x35, 0xfa, 0x00, 0x20, 0xd9, 0x24, 0xd7, 0x9f, 0xf2, 0x21, 0x52, 0xc4, 0x39, 0x84, 0x5d, 0xc0,
	0x74, 0x45, 0x7d, 0x83, 0x5f, 0x31, 0x31, 0xbe, 0x80, 0x4b, 0xb8, 0xcf, 0x2e, 0xdb, 0xcb, 0xfc,
	0x69, 0x49, 0x87, 0xd2, 0x1e, 0xbf, 0x77, 0x95, 0x8c, 0xd0, 0xf9, 0x74, 0xaa, 0xff, 0x2d, 0x80,
	0xb4, 0x7e, 0x5c, 0xd0, 0x29, 0x9c, 0x2c, 0xf3, 0x1a, 0xb9, 0xc9, 0x21, 0x70, 0x93, 0xe3, 0x25,
	0xd9, 0xcd, 0x66, 0xc8, 0xc7, 0xac, 0x11, 0x59, 0x4c, 0x9a, 0x35, 0x1e, 0x94, 0xd2, 0x92, 0x89,
	0xd3, 0xae, 0xf4, 0x61, 0x77, 0xad, 0x0f, 0x1d, 0xc8, 0x25, 0x30, 0x32, 0x99, 0xc8, 0x65, 0xb9,
	0x24, 0xa3, 0x84, 0xa9, 0xff, 0x00, 0xb5, 0x6d, 0x27, 0x93, 0x4d, 0x6d, 0xea, 0x39, 0x4b, 0x07,
	0x81, 0x3b, 0x94, 0xa9, 0xe7, 0xa4, 0xa1, 0x4c, 0xe2, 0x93, 0x77, 0xc6, 0xda, 0x9e, 0x94, 0x7d,
	0xf2, 0x2e, 0x73, 0xff, 0x63, 0x07, 0xa4, 0xf5, 0xc3, 0x89, 0x7a, 0x50, 0x08, 0x88, 0x19, 0x52,
	0x3f, 0xf9, 0x92, 0xbc, 0xfc, 0x0f, 0x27, 0xba, 0x8d, 0x79, 0x04, 0x4e, 0x22, 0xd9, 0x50, 0xa6,
	0xd1, 0x15, 0x09, 0x78, 0xb3, 0x92, 0x4e, 0x95, 0x38, 0xc2, 0x3c, 0xd0, 0xff, 0x20, 0x5e, 0x18,
	0x16, 0x8d, 0x92, 0x4f, 0x4a, 0x91, 0x03, 0x3d, 0x1a, 0x35, 0x7f, 0x17, 0xa0, 0x10, 0xdb, 0xa1,
	0x0a, 0x94, 0x27, 0xda, 0x78, 0xa4, 0xf4, 0xd5, 0xd7, 0xaa, 0x32, 0x90, 0x1e, 0xa0, 0x03, 0x28,
	0xf6, 0x86, 0xba, 0xa1, 0x0e, 0xce, 0x15, 0x49, 0x40, 0x35, 0x90, 0xd8, 0x93, 0xa1, 0x77, 0xc7,
	0x5f, 0x19, 0xdd, 0xde, 0x58, 0xd1, 0x74, 0x69, 0x07, 0x3d, 0x82, 0x63, 0x3c, 0xd1, 0x34, 0x55,
	0x3b, 0x5b, 0x21, 0x76, 0xd1, 0x13, 0x78, 0xac, 0x6a, 0xfd, 0xa1, 0x36, 0x56, 0xc7, 0xba, 0xa2,
	0xe9, 0x06, 0x73, 0x7a, 0x3d, 0xc4, 0x5c, 0x25, 0x89, 0x1b, 0x34, 0x0f, 0x66, 0x7c, 0x6f, 0xa8,
	0x4b, 0x7b, 0xf5, 0x63, 0xa8, 0x6e, 0x5c, 0xb9, 0xe6, 0xf7, 0x70, 0xb8, 0x32, 0x9a, 0x91, 0x0c,
	0xb5, 0x7e, 0x57, 0x57, 0xce, 0x86, 0xf8, 0x5b, 0x63, 0xb5, 0xf4, 0x87, 0x50, 0xcd, 0x98, 0xf1,
	0x37, 0x5d, 0x7c, 0xa1, 0x6a, 0x67, 0xd2, 0xfb, 0x1d, 0x24, 0xc3, 0x71, 0x86, 0xbf, 0x1d, 0xf7,
	0xdf, 0x28, 0x83, 0xc9, 0xb9, 0x82, 0xa5, 0xf7, 0x62, 0xf3, 0x2f, 0x01, 0x4a, 0xd9, 0x27, 0x7b,
	0xb3, 0x17, 0x0f, 0xa1, 0x9a, 0xfa, 0x18, 0x8a, 0xf6, 0x76, 0xa2, 0x4c, 0x94, 0x81, 0xf4, 0x33,
	0x7b, 0x7f, 0x94, 0xe1, 0xfd, 0xe1, 0xc5, 0xe8, 0x5c, 0xd1, 0x95, 0x81, 0xf4, 0x0b, 0xcf, 0xb4,
	0x4c, 0x60, 0x74, 0xc7, 0x63, 0xf5, 0x4c, 0x63, 0x21, 0x22, 0x7a, 0x0c, 0xb5, 0x1c, 0x33, 0xc2,
	0x8a, 0x72, 0x31, 0xe2, 0x41, 0xac, 0x2b, 0x72, 0x8e, 0xc2, 0xca, 0x08, 0xab, 0x43, 0xac, 0xea,
	0xea, 0x77, 0xca, 0x40, 0xfa, 0x55, 0x44, 0x75, 0x38, 0xc9, 0xd1, 0x13, 0x2d, 0x73, 0xfd, 0x4d,
	0xec, 0x95, 0x60, 0x3f, 0x19, 0x63, 0x5f, 0x8a, 0xc5, 0x8a, 0xe4, 0x34, 0x3f, 0x4b, 0xde, 0xe7,
	0xdc, 0xe5, 0x37, 0xb9, 0xc0, 0x7f, 0x84, 0xc4, 0x73, 0xa8, 0x7c, 0x8a, 0x36, 0x8f, 0x19, 0x4e,
	0x14, 0x56, 0x81, 0xff, 0xa2, 0x78, 0xf5, 0x6f, 0x00, 0x00, 0x00, 0xff, 0xff, 0xda, 0xad, 0xb8,
	0xac, 0x56, 0x09, 0x00, 0x00,
}
