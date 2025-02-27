// Copyright 2017 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package structs

import (
	"time"
)

// LFSLock represent a lock
// for use with the locks API.
//
//swagger:model
type LFSLock struct {
	ID       string        `json:"id"`
	Path     string        `json:"path"`
	LockedAt time.Time     `json:"locked_at"`
	Owner    *LFSLockOwner `json:"owner"`
}

// LFSLockOwner represent a lock owner
// for use with the locks API.
type LFSLockOwner struct {
	Name string `json:"name"`
}

// LFSLockRequest contains the path of the lock to create
// https://github.com/git-lfs/git-lfs/blob/master/docs/api/locking.md#create-lock
type LFSLockRequest struct {
	Path string `json:"path"`
}

// LFSLockResponse represent a lock created
// https://github.com/git-lfs/git-lfs/blob/master/docs/api/locking.md#create-lock
type LFSLockResponse struct {
	Lock *LFSLock `json:"lock"`
}

// LFSLockList represent a list of lock requested
// https://github.com/git-lfs/git-lfs/blob/master/docs/api/locking.md#list-locks
// swagger:response LFSLockList
type LFSLockList struct {
	Locks []*LFSLock `json:"locks"`
	Next  string     `json:"next_cursor,omitempty"`
}

// LFSLockListVerify represent a list of lock verification requested
// https://github.com/git-lfs/git-lfs/blob/master/docs/api/locking.md#list-locks-for-verification
type LFSLockListVerify struct {
	Ours   []*LFSLock `json:"ours"`
	Theirs []*LFSLock `json:"theirs"`
	Next   string     `json:"next_cursor,omitempty"`
}

// LFSLockError contains information on the error that occurs
type LFSLockError struct {
	Message       string   `json:"message"`
	Lock          *LFSLock `json:"lock,omitempty"`
	Documentation string   `json:"documentation_url,omitempty"`
	RequestID     string   `json:"request_id,omitempty"`
}

// LFSLockDeleteRequest contains params of a delete request
// https://github.com/git-lfs/git-lfs/blob/master/docs/api/locking.md#delete-lock
type LFSLockDeleteRequest struct {
	Force bool `json:"force"`
}

// UnlockList list of LFS lock ids to unlock
type UnlockList struct {
	// ids to unlock
	//
	// in: body
	LockIds []int64 `json:"lock_ids"`
	Force   bool    `json:"force"`
}

// LFSUnlockedList gives the locks that were actually unlocked
// swagger:response LFSUnlockedList
type LFSUnlockedList struct {
	Locks []*LFSLock `json:"locks"`
}
