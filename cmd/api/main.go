package main

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// SmartContract provides functions for managing the supply chain
type SmartContract struct {
	contractapi.Contract
}

// Product describes basic details of what makes up a product
type Product struct {
	ID          string `json:"id"`
	Description string `json:"description"`
	Owner       string `json:"owner"`
	Status      string `json:"status"`
}

// InitLedger adds a base set of products to the ledger
func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	products := []Product{
		{ID: "product1", Description: "Laptop", Owner: "Manufacturer", Status: "Manufactured"},
		{ID: "product2", Description: "Phone", Owner: "Manufacturer", Status: "Manufactured"},
	}

	for _, product := range products {
		productJSON, err := json.Marshal(product)
		if err != nil {
			return err
		}

		err = ctx.GetStub().PutState(product.ID, productJSON)
		if err != nil {
			return fmt.Errorf("Failed to put to world state. %s", err.Error())
		}
	}

	return nil
}

// CreateProduct adds a new product to the world state
func (s *SmartContract) CreateProduct(ctx contractapi.TransactionContextInterface, id string, description string) error {
	exists, err := s.ProductExists(ctx, id)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("Product %s already exists", id)
	}

	product := Product{
		ID:          id,
		Description: description,
		Owner:       "Manufacturer",
		Status:      "Manufactured",
	}
	productJSON, err := json.Marshal(product)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, productJSON)
}

// TransferProduct updates the owner and status of a product
func (s *SmartContract) TransferProduct(ctx contractapi.TransactionContextInterface, id string, newOwner string, newStatus string) error {
	product, err := s.ReadProduct(ctx, id)
	if err != nil {
		return err
	}

	product.Owner = newOwner
	product.Status = newStatus

	productJSON, err := json.Marshal(product)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, productJSON)
}

// ReadProduct returns the product stored in the world state with given id
func (s *SmartContract) ReadProduct(ctx contractapi.TransactionContextInterface, id string) (*Product, error) {
	productJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return nil, fmt.Errorf("Failed to read from world state. %s", err.Error())
	}
	if productJSON == nil {
		return nil, fmt.Errorf("Product %s does not exist", id)
	}

	var product Product
	err = json.Unmarshal(productJSON, &product)
	if err != nil {
		return nil, err
	}

	return &product, nil
}

// ProductExists returns true when product with given ID exists in world state
func (s *SmartContract) ProductExists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
	productJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return false, fmt.Errorf("Failed to read from world state. %s", err.Error())
	}

	return productJSON != nil, nil
}

func main() {
	chaincode, err := contractapi.NewChaincode(new(SmartContract))
	if err != nil {
		fmt.Printf("Error creating supplychain chaincode: %s", err.Error())
		return
	}

	if err := chaincode.Start(); err != nil {
		fmt.Printf("Error starting supplychain chaincode: %s", err.Error())
	}
}
