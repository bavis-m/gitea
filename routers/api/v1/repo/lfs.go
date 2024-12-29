package repo

import (
	"net/http"
	"strconv"

	git_model "code.gitea.io/gitea/models/git"
	"code.gitea.io/gitea/modules/log"
	api "code.gitea.io/gitea/modules/structs"
	"code.gitea.io/gitea/modules/web"
	"code.gitea.io/gitea/services/context"
	"code.gitea.io/gitea/services/convert"
)

// ListLFSLocks list all the locked files in a repository
func ListLFSLocks(ctx *context.APIContext) {
	// swagger:operation GET /repos/{owner}/{repo}/lfs/locks repository repoListLFSLocks
	// ---
	// summary: List a repository's LFS locks
	// produces:
	// - application/json
	// parameters:
	// - name: owner
	//   in: path
	//   description: owner of the repo
	//   type: string
	//   required: true
	// - name: repo
	//   in: path
	//   description: name of the repo
	//   type: string
	//   required: true
	// - name: cursor
	//   in: query
	//   description: cursor to fetch from
	//   type: integer
	//   required: false
	// - name: limit
	//   in: query
	//   description: how many locks to fetch
	//   type: integer
	//   required: false
	// responses:
	//   "200":
	//     "$ref": "#/responses/LFSLockList"
	//   "404":
	//     "$ref": "#/responses/notFound"

	cursor := ctx.FormInt("cursor")
	limit := ctx.FormInt("limit")

	lockList, err := git_model.GetLFSLockByRepoID(ctx, ctx.Repo.Repository.ID, cursor, limit)
	if err != nil {
		log.Error("Unable to list locks for repository ID[%d]: Error: %v", ctx.Repo.Repository.ID, err)
		ctx.JSON(http.StatusInternalServerError, api.LFSLockError{
			Message: "unable to list locks : Internal Server Error",
		})
		return
	}

	lockListAPI := make([]*api.LFSLock, len(lockList))
	next := ""
	for i, l := range lockList {
		lockListAPI[i] = convert.ToLFSLock(ctx, l)
	}
	if limit > 0 && len(lockList) == limit {
		next = strconv.Itoa(cursor + 1)
	}
	ctx.JSON(http.StatusOK, api.LFSLockList{
		Locks: lockListAPI,
		Next:  next,
	})
}

// UnlockLFSLocks unlocks multiple LFS locks at once
func UnlockLFSLocks(ctx *context.APIContext) {
	// swagger:operation POST /repos/{owner}/{repo}/lfs/unlock_locks repository repoUnlockLFSLocks
	// ---
	// summary: UnlockLFSLocks unlocks multiple LFS locks at once
	// produces:
	// - application/json
	// parameters:
	// - name: owner
	//   in: path
	//   description: owner of the repo
	//   type: string
	//   required: true
	// - name: repo
	//   in: path
	//   description: name of the repo
	//   type: string
	//   required: true
	// - name: body
	//   in: body
	//   schema:
	//     "$ref": "#/definitions/UnlockList"
	// responses:
	//   "200":
	//     "$ref": "#/responses/LFSUnlockedList"
	//   "404":
	//     "$ref": "#/responses/notFound"

	lockIds, ok := web.GetForm(ctx).(*api.UnlockList)
	if !ok {
		log.Error("Invalid data")
		ctx.JSON(http.StatusBadRequest, api.LFSLockError{
			Message: "Unable to unlock : Invalid Data",
		})
		return
	}

	lockListAPI := make([]*api.LFSLock, 0, len(lockIds.LockIds))
	for _, id := range lockIds.LockIds {
		lock, err := git_model.DeleteLFSLockByID(ctx, id, ctx.Repo.Repository, ctx.Doer /*req.Force*/, false)
		if git_model.IsErrLFSLockNotExist(err) {
			// this is OK
		} else if err != nil {
			log.Error("Unable to unlock lock ID[%d]: Error: %v", id, err)
			ctx.JSON(http.StatusInternalServerError, api.LFSLockError{
				Message: "Unable to unlock : Internal Server Error",
			})
			return
		} else {
			lockListAPI = append(lockListAPI, convert.ToLFSLock(ctx, lock))
		}
	}

	ctx.JSON(http.StatusOK, api.LFSUnlockedList{
		Locks: lockListAPI,
	})
}

// LFSLockList
// swagger:response LFSLockList
type swaggerLFSLockList struct {
	// in: body
	Body api.LFSLockList `json:"body"`
}

// LFSUnlockedList
// swagger:response LFSUnlockedList
type swaggerLFSUnlockedList struct {
	// in: body
	Body api.LFSUnlockedList `json:"body"`
}
