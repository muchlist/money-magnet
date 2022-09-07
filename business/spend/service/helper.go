package service

import (
	"strings"

	"github.com/google/uuid"
)

// isCanEditOrWatch read pocketRoles and validate those roles can edit or just watch
// isCanEditOrWatch can be move to claims function if needed globally
func isCanEditOrWatch(pocketID uuid.UUID, pocketRoles []string) (edit bool, watch bool) {
	pocketStr := pocketID.String()
	canEdit := false
	canWatch := false

	for _, role := range pocketRoles {
		pocketIDandAccess := strings.Split(role, ":")
		if len(pocketIDandAccess) != 2 {
			return false, false
		}

		id := pocketIDandAccess[0]
		access := pocketIDandAccess[1]

		if id == pocketStr {
			if access == "edit" {
				canEdit = true
			}
			if access == "watch" {
				canWatch = true
			}
		}
	}
	return canEdit, canWatch
}
