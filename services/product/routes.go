package product

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/xhoang0509/ecom-api/services/auth"
	"github.com/xhoang0509/ecom-api/types"
	"github.com/xhoang0509/ecom-api/utils"
	"net/http"
	"strconv"
)

type Handler struct {
	store     types.ProductStore
	userStore types.UserStore
}

func NewHandler(store types.ProductStore, userStore types.UserStore) *Handler {
	return &Handler{
		store:     store,
		userStore: userStore,
	}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/products", h.handleGetProducts).Methods(http.MethodGet)
	router.HandleFunc("/products/{productID}", h.handleGetProduct).Methods(http.MethodGet)

	router.HandleFunc("/products", auth.WithJWTAuth(h.handleCreateProduct, h.userStore)).Methods(http.MethodPost)
}

func (h *Handler) handleGetProducts(w http.ResponseWriter, r *http.Request) {
	products, err := h.store.GetProducts()
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJson(w, http.StatusOK, products)
}

func (h *Handler) handleGetProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	str, ok := vars["productID"]

	if !ok {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("missing product ID"))
		return
	}

	productID, err := strconv.Atoi(str)

	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid product ID"))
		return
	}

	product, err := h.store.GetProductByID(productID)
	if err != nil {
		utils.WriteError(w, http.StatusNotFound, fmt.Errorf("not found product with ID %d", productID))
		return
	}

	utils.WriteJson(w, http.StatusOK, product)
}

func (h *Handler) handleCreateProduct(w http.ResponseWriter, r *http.Request) {
	//	Get JSON Payload
	var product types.CreateProductPayload
	if err := utils.ParseJson(r, &product); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	//	validate payload
	if err := utils.Validate.Struct(product); err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload: %s", err))
		return
	}

	err := h.store.CreateProduct(types.CreateProductPayload{
		Name:        product.Name,
		Description: product.Description,
		Image:       product.Image,
		Price:       product.Price,
		Quantity:    product.Quantity,
	})
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("create product failed: %v", err))
		return
	}

	utils.WriteJson(w, http.StatusCreated, map[string]string{"message": "create product success!"})
}
