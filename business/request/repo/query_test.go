package repo

import (
	"testing"

	sq "github.com/Masterminds/squirrel"
	"github.com/go-playground/assert/v2"
	"github.com/muchlist/moneymagnet/business/request/model"
)

// go test -v -timeout 30s -run ^TestComplexQuery$ github.com/muchlist/moneymagnet/business/request/repo
func TestComplexQuery(t *testing.T) {
	sb := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	findBy := model.FindBy{
		PocketIDs:   []string{"1", "2"},
		ApproverID:  "approverIDExample",
		RequesterID: "",
	}
	expected := `SELECT count(*) OVER(), id, requester_id, approver_id, pocket_id, pocket_name, is_approved, is_rejected, created_at, updated_at FROM requests WHERE (approver_id = $1 OR pocket_id IN ($2,$3))`

	// where builder
	var orBuilder sq.Or
	andBuilder := sq.Eq{} // NOT used in this test
	var orValidCount int
	if findBy.ApproverID != "" {
		orBuilder = append(orBuilder, sq.Eq{keyApproverID: findBy.ApproverID})
		andBuilder[keyApproverID] = findBy.ApproverID
		orValidCount++
	}
	if len(findBy.PocketIDs) != 0 {
		orBuilder = append(orBuilder, sq.Eq{keyPocketID: findBy.PocketIDs})
		andBuilder[keyPocketID] = findBy.PocketIDs
		orValidCount++
	}
	if findBy.RequesterID != "" {
		orBuilder = append(orBuilder, sq.Eq{keyRequesterID: findBy.RequesterID})
		andBuilder[keyRequesterID] = findBy.RequesterID
		orValidCount++
	}

	// IF use or but input less than 2 return not vali

	sqlStatement, args, err := sb.Select(
		"count(*) OVER()",
		keyID,
		keyRequesterID,
		keyApproverID,
		keyPocketID,
		keyPocketName,
		keyIsApproved,
		keyIsRejected,
		keyCreatedAt,
		keyUpdatedAt,
	).
		From(keyTable).
		Where(orBuilder).
		ToSql()

	if err != nil {
		t.Errorf("build query find request: %s", err)
		return
	}

	t.Log(sqlStatement)
	t.Log(args)

	assert.Equal(t, expected, sqlStatement)
}
