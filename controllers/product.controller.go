package controllers

// add image url validator middleware
// add a product search controller

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/harshvardha/e-commerce/internal/database"
	"github.com/harshvardha/e-commerce/utility"
)

func (apiConfig *ApiConfig) HandleListProduct(w http.ResponseWriter, r *http.Request, IDs ID, newAccessToken string) {
	// decoding the request body
	decoder := json.NewDecoder(r.Body)
	params := CreateOrUpdateProductRequest{}
	err := decoder.Decode(&params)
	if err != nil {
		utility.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	// listing a new product
	if len(params.Description) == 0 || len(params.ImageUrls) == 0 || len(params.Characteristics) == 0 {
		utility.RespondWithError(w, http.StatusBadRequest, "Product Description Or Image Urls Or Characteristics Cannot Be Empty")
		return
	}
	productDescription, err := json.Marshal(params.Description)
	if err != nil {
		utility.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	imageUrls, err := json.Marshal(params.ImageUrls)
	if err != nil {
		utility.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	newProduct, err := apiConfig.DB.ListProduct(r.Context(), database.ListProductParams{
		Name:        params.Name,
		Description: productDescription,
		Price:       params.Price,
		ImageUrls:   imageUrls,
		StockAmount: params.StockAmount,
		StoreID:     IDs.StoreID,
		CategoryID:  params.CategoryID,
	})

	if err != nil {
		utility.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	// registering product characteristics
	characteristics, err := json.Marshal(params.Characteristics)
	if err != nil {
		utility.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	productCharacteristics, err := apiConfig.DB.CreateCharacteristics(r.Context(), database.CreateCharacteristicsParams{
		Description: characteristics,
		ProductID:   newProduct.ID,
	})
	if err != nil {
		utility.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utility.RespondWithJson(w, http.StatusCreated, ProductResponse{
		ID:              newProduct.ID,
		Name:            newProduct.Name,
		Description:     newProduct.Description,
		Characteristics: productCharacteristics.Description,
		Price:           newProduct.Price,
		ImageUrls:       newProduct.ImageUrls,
		StockAmount:     newProduct.StockAmount,
		StoreID:         newProduct.StoreID,
		CategoryID:      newProduct.CategoryID,
		AccessToken:     newAccessToken,
		CreatedAt:       newProduct.CreatedAt,
		UpdatedAt:       newProduct.UpdatedAt,
	})
}

func (apiConfig *ApiConfig) HandleUpdateProduct(w http.ResponseWriter, r *http.Request, IDs ID, newAccessToken string) {
	// fetching the product id to update
	productIDString := r.PathValue("product_id")
	productCharacteristicsIDString := r.PathValue("product_characteristic_id")
	if productIDString == "" || productCharacteristicsIDString == "" {
		utility.RespondWithError(w, http.StatusBadRequest, "Invalid product id or characteristic id")
		return
	}
	productID, err := uuid.Parse(productIDString)
	if err != nil {
		utility.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	productCharacteristicsID, err := uuid.Parse(productCharacteristicsIDString)
	if err != nil {
		utility.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	// decoding request body
	decoder := json.NewDecoder(r.Body)
	params := CreateOrUpdateProductRequest{}
	err = decoder.Decode(&params)
	if err != nil {
		utility.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	// checking which fields to update for product
	existingProduct, err := apiConfig.DB.GetProductById(r.Context(), productID)
	if err != nil {
		utility.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	updateProduct := database.UpdateProductParams{
		Name:        existingProduct.Name,
		Description: existingProduct.Description,
		Price:       existingProduct.Price,
		StockAmount: existingProduct.StockAmount,
		ImageUrls:   existingProduct.ImageUrls,
		CategoryID:  existingProduct.CategoryID,
		ID:          productID,
	}

	existingProductCharacteristics, err := apiConfig.DB.GetProductCharacteristic(r.Context(), database.GetProductCharacteristicParams{
		ID:        productCharacteristicsID,
		ProductID: productID,
	})
	if err != nil {
		utility.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	updateProductCharacteristics := database.UpdateCharacteristicsParams{
		Description: existingProductCharacteristics,
		ID:          productCharacteristicsID,
		ProductID:   productID,
	}

	if params.Name != "" {
		updateProduct.Name = params.Name
	}
	if len(params.Description) > 0 {
		newProductDescription, err := json.Marshal(params.Description)
		if err != nil {
			utility.RespondWithError(w, http.StatusBadRequest, "Invalid Product Description")
			return
		}
		updateProduct.Description = newProductDescription
	}
	if len(params.Characteristics) > 0 {
		newProductCharacteristic, err := json.Marshal(params.Characteristics)
		if err != nil {
			utility.RespondWithError(w, http.StatusBadRequest, "Invalid Product Characteristic")
			return
		}
		updateProductCharacteristics.Description = newProductCharacteristic
	}
	if params.Price > 0 {
		updateProduct.Price = params.Price
	}
	if len(params.ImageUrls) > 0 {
		newProductImageUrls, err := json.Marshal(params.ImageUrls)
		if err != nil {
			utility.RespondWithError(w, http.StatusBadGateway, "Invalid Image Urls")
			return
		}
		updateProduct.ImageUrls = newProductImageUrls
	}
	if params.StockAmount >= 0 {
		updateProduct.StockAmount = params.StockAmount
	}
	if params.CategoryID != uuid.Nil {
		updateProduct.CategoryID = params.CategoryID
	}

	// updating product
	updatedProduct, err := apiConfig.DB.UpdateProduct(r.Context(), updateProduct)
	if err != nil {
		utility.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// updating product characteristic
	updatedProductCharacteristic, err := apiConfig.DB.UpdateCharacteristics(r.Context(), updateProductCharacteristics)
	if err != nil {
		utility.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utility.RespondWithJson(w, http.StatusOK, ProductResponse{
		ID:              updatedProduct.ID,
		Name:            updatedProduct.Name,
		Description:     updateProduct.Description,
		Characteristics: updatedProductCharacteristic,
		Price:           updatedProduct.Price,
		ImageUrls:       updatedProduct.ImageUrls,
		StockAmount:     updatedProduct.StockAmount,
		StoreID:         updatedProduct.StoreID,
		CategoryID:      updatedProduct.CategoryID,
		AccessToken:     newAccessToken,
		CreatedAt:       updatedProduct.CreatedAt,
		UpdatedAt:       updatedProduct.UpdatedAt,
	})
}

func (apiConfig *ApiConfig) HandleRemoveProduct(w http.ResponseWriter, r *http.Request, IDs ID, newAccessToken string) {
	// fetching the product id
	productIDString := r.PathValue("product_id")
	if productIDString == "" {
		utility.RespondWithError(w, http.StatusBadRequest, "Invalid product id")
		return
	}
	productID, err := uuid.Parse(productIDString)
	if err != nil {
		utility.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// removing the product from store
	err = apiConfig.DB.RemoveProduct(r.Context(), database.RemoveProductParams{
		ID:      productID,
		StoreID: IDs.StoreID,
	})
	if err != nil {
		utility.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	utility.RespondWithJson(w, http.StatusOK, EmptyResponse{
		AccessToken: newAccessToken,
	})
}

func (apiConfig *ApiConfig) HandleGetProductInformation(w http.ResponseWriter, r *http.Request, IDs ID, newAccessToken string) {
	// fetching the product id
	productIDString := r.PathValue("product_id")

	if productIDString == "" {
		utility.RespondWithError(w, http.StatusBadRequest, "Invalid product id")
		return
	}
	productID, err := uuid.Parse(productIDString)
	if err != nil {
		utility.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// fetching all the information about the product
	product, err := apiConfig.DB.GetProductById(r.Context(), productID)
	if err != nil {
		utility.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	// fetching all the product characteristics
	productCharacteristics, err := apiConfig.DB.GetAllProductCharacteristics(r.Context(), productID)
	if err != nil {
		utility.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	productResponse := struct {
		Product                database.Product          `json:"product"`
		ProductCharacteristics []database.Characteristic `json:"product_characteristics"`
	}{
		Product:                product,
		ProductCharacteristics: productCharacteristics,
	}

	utility.RespondWithJson(w, http.StatusOK, productResponse)
}

func (apiConfig *ApiConfig) HandleGetProductsByCategory(w http.ResponseWriter, r *http.Request, IDs ID, newAccessToken string) {
	// fetching the category id
	categoryIDString := r.PathValue("category_id")
	if categoryIDString == "" {
		utility.RespondWithError(w, http.StatusBadRequest, "Invalid category id")
		return
	}
	categoryID, err := uuid.Parse(categoryIDString)
	if err != nil {
		utility.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	// fetching the products
	products, err := apiConfig.DB.GetProductsByCategory(r.Context(), categoryID)
	if err != nil {
		utility.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	utility.RespondWithJson(w, http.StatusOK, ProductsByCategoryOrStoreID{
		Products:    products,
		AccessToken: newAccessToken,
	})
}

func (apiConfig *ApiConfig) HandleGetProductsByStoreID(w http.ResponseWriter, r *http.Request, IDs ID, newAccessToken string) {
	if IDs.StoreID == uuid.Nil {
		utility.RespondWithError(w, http.StatusBadRequest, "Invalid store id")
		return
	}

	// fetching the products
	products, err := apiConfig.DB.GetProductsByStoreId(r.Context(), IDs.StoreID)
	if err != nil {
		utility.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	utility.RespondWithJson(w, http.StatusOK, ProductsByCategoryOrStoreID{
		Products:    products,
		AccessToken: newAccessToken,
	})
}
