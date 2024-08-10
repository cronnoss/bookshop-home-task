package httpserver

import (
	"encoding/json"
	"net/http"

	"github.com/cronnoss/bookshop-home-task/internal/app/common/server"
)

// @Summary UpdateCart
// @Security ApiKeyAuth
// @Tags cart
// @Description update cart
// @ID update-cart
// @Accept  json
// @Produce  json
// @Param input body CartRequest true "cart info"
// @Success 200 {object} CartResponse
// @Failure 400,404 {object} server.ErrorResponse
// @Failure 401 {object} server.ErrorResponse
// @Failure 500 {object} server.ErrorResponse
// @Router /cart [post]
func (h HTTPServer) UpdateCart(w http.ResponseWriter, r *http.Request) {
	user, err := getUserFromContext(r.Context())
	if err != nil {
		server.BadRequest("invalid-user", err, w, r)
		return
	}

	var cartRequest CartRequest
	if err := json.NewDecoder(r.Body).Decode(&cartRequest); err != nil {
		server.BadRequest("invalid-json", err, w, r)
		return
	}

	_, err = h.userService.GetUserByID(r.Context(), user.ID)
	if err != nil {
		server.RespondWithError(err, w, r)
		return
	}

	cart, err := toDomainCart(user.ID, cartRequest)
	if err != nil {
		server.RespondWithError(err, w, r)
		return
	}

	updatedCart, err := h.cartService.UpdateCartAndStocks(r.Context(), cart)
	if err != nil {
		server.RespondWithError(err, w, r)
		return
	}

	response := toResponseCart(updatedCart)

	server.RespondOK(response, w, r)
}

// @Summary Checkout
// @Security ApiKeyAuth
// @Tags cart
// @Description checkout
// @ID checkout
// @Accept  json
// @Produce  json
// @Success 200 {object} map[string]bool
// @Failure 400,404 {object} server.ErrorResponse
// @Failure 401 {object} server.ErrorResponse
// @Failure 500 {object} server.ErrorResponse
// @Router /checkout [post]
func (h HTTPServer) Checkout(w http.ResponseWriter, r *http.Request) {
	user, err := getUserFromContext(r.Context())
	if err != nil {
		server.BadRequest("invalid-user", err, w, r)
		return
	}

	err = h.cartService.Checkout(r.Context(), user.ID)
	if err != nil {
		server.RespondWithError(err, w, r)
		return
	}

	server.RespondOK(map[string]bool{"ok": true}, w, r)
}
