package service

import (
	"net/http"

	"github.com/ArkaniLoveCoding/Shcool-manajement/types"
	"github.com/ArkaniLoveCoding/Shcool-manajement/utils"
)

// this is for router that token is not verified in their function!

type HandleRequest struct {
	db types.UserStore
}

func NewHandlerUser(db types.UserStore) *HandleRequest {
	return &HandleRequest{db: db}
}

// this is for router that token is verified in their function!

type HandleRequestForAuthenticate struct {
	db types.UserStore
}

func NewHandlerUserForAuthenticate (db types.UserStore) *HandleRequestForAuthenticate {
	return &HandleRequestForAuthenticate{db: db}
}

func (h *HandleRequest) Register_Bp(w http.ResponseWriter, r *http.Request) {

	//create the payload of the json
	var payload types.Register
	if err := utils.DecodeData(r, &payload); err != nil {
		
	}

}

