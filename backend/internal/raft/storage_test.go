package raft

import (
	"testing"
)

func TestMemoryStorage_Term(t *testing.T) {
	s := NewMemoryStorage()

	// 测试初始任期
	term, err := s.GetCurrentTerm()
	if err != nil {
		t.Fatalf("获取初始任期失败: %v", err)
	}
	if term != 0 {
		t.Errorf("期望初始任期为 0，实际为 %d", term)
	}

	// 测试设置任期
	if err := s.SetCurrentTerm(5); err != nil {
		t.Fatalf("设置任期失败: %v", err)
	}

	term, err = s.GetCurrentTerm()
	if err != nil {
		t.Fatalf("获取任期失败: %v", err)
	}
	if term != 5 {
		t.Errorf("期望任期为 5，实际为 %d", term)
	}
}

func TestMemoryStorage_VotedFor(t *testing.T) {
	s := NewMemoryStorage()

	// 测试初始投票
	votedFor, err := s.GetVotedFor()
	if err != nil {
		t.Fatalf("获取初始投票失败: %v", err)
	}
	if votedFor != "" {
		t.Errorf("期望初始投票为空，实际为 %s", votedFor)
	}

	// 测试设置投票
	if err := s.SetVotedFor("node1"); err != nil {
		t.Fatalf("设置投票失败: %v", err)
	}

	votedFor, err = s.GetVotedFor()
	if err != nil {
		t.Fatalf("获取投票失败: %v", err)
	}
	if votedFor != "node1" {
		t.Errorf("期望投票为 node1，实际为 %s", votedFor)
	}
}

func TestMemoryStorage_Logs(t *testing.T) {
	s := NewMemoryStorage()

	// 测试空日志
	first, err := s.FirstIndex()
	if err != nil {
		t.Fatalf("获取第一个日志索引失败: %v", err)
	}
	if first != 0 {
		t.Errorf("期望第一个日志索引为 0，实际为 %d", first)
	}

	last, err := s.LastIndex()
	if err != nil {
		t.Fatalf("获取最后一个日志索引失败: %v", err)
	}
	if last != 0 {
		t.Errorf("期望最后一个日志索引为 0，实际为 %d", last)
	}

	// 测试存储日志
	logs := []*LogEntry{
		{Term: 1, Index: 1, Command: "cmd1"},
		{Term: 1, Index: 2, Command: "cmd2"},
		{Term: 2, Index: 3, Command: "cmd3"},
	}
	if err := s.StoreLogs(logs); err != nil {
		t.Fatalf("存储日志失败: %v", err)
	}

	// 验证日志索引
	first, _ = s.FirstIndex()
	last, _ = s.LastIndex()
	if first != 1 {
		t.Errorf("期望第一个日志索引为 1，实际为 %d", first)
	}
	if last != 3 {
		t.Errorf("期望最后一个日志索引为 3，实际为 %d", last)
	}

	// 测试获取日志
	log, err := s.GetLog(2)
	if err != nil {
		t.Fatalf("获取日志失败: %v", err)
	}
	if log.Term != 1 || log.Index != 2 || log.Command != "cmd2" {
		t.Errorf("日志内容不匹配，期望 {1,2,cmd2}，实际为 %v", log)
	}

	// 测试删除日志范围
	if err := s.DeleteRange(2, 3); err != nil {
		t.Fatalf("删除日志范围失败: %v", err)
	}

	last, _ = s.LastIndex()
	if last != 1 {
		t.Errorf("删除日志后，期望最后一个日志索引为 1，实际为 %d", last)
	}
}

func TestMemoryStorage_Snapshot(t *testing.T) {
	s := NewMemoryStorage()

	// 测试空快照
	_, err := s.GetSnapshot()
	if err != ErrNoSnapshot {
		t.Errorf("期望错误为 ErrNoSnapshot，实际为 %v", err)
	}

	// 测试存储快照
	snapshot := []byte("snapshot data")
	if err := s.StoreSnapshot(snapshot); err != nil {
		t.Fatalf("存储快照失败: %v", err)
	}

	// 测试获取快照
	data, err := s.GetSnapshot()
	if err != nil {
		t.Fatalf("获取快照失败: %v", err)
	}
	if string(data) != "snapshot data" {
		t.Errorf("快照数据不匹配，期望 'snapshot data'，实际为 %s", string(data))
	}
}
