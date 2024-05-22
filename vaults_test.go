package dvls

import (
	"reflect"
	"testing"
)

const testNewVaultId string = "eabd3646-acf8-44a4-9ba0-991df147c209"

var testNewVaultPassword string = "5w:mr6kPj"

var testVault Vault = Vault{
	Name:        "go-dvls tests",
	Description: "Test Vault",
}

var testNewVault Vault = Vault{
	ID:          testNewVaultId,
	Name:        "go-dvls tests new",
	Description: "Test",
}

func Test_Vaults(t *testing.T) {
	testVault.ID = testVaultId
	t.Run("GetVault", test_GetVault)
	t.Run("NewVault", test_NewVault)
	t.Run("UpdateVault", test_UpdateVault)
	t.Run("DeleteVault", test_DeleteVault)
}

func test_GetVault(t *testing.T) {
	vault, err := testClient.GetVault(testVaultId)
	if err != nil {
		t.Fatal(err)
	}

	testVault.CreationDate = vault.CreationDate
	testVault.ModifiedDate = vault.ModifiedDate

	if !reflect.DeepEqual(testVault, vault) {
		t.Fatalf("fetched vault did not match test vault. Expected %#v, got %#v", testVault, vault)
	}
}

func test_NewVault(t *testing.T) {
	err := testClient.NewVault(testNewVault, nil)
	if err != nil {
		t.Fatal(err)
	}

	vault, err := testClient.GetVault(testNewVault.ID)
	if err != nil {
		t.Fatal(err)
	}

	vault.CreationDate = testNewVault.CreationDate
	vault.ModifiedDate = testNewVault.ModifiedDate

	if !reflect.DeepEqual(testNewVault, vault) {
		t.Fatalf("fetched vault did not match test vault. Expected %#v, got %#v", testNewVault, vault)
	}
}

func test_UpdateVault(t *testing.T) {
	testNewVault.Name = "go-dvls tests new updated"
	testNewVault.Description = "Test updated"
	options := VaultOptions{Password: &testNewVaultPassword}

	err := testClient.UpdateVault(testNewVault, &options)
	if err != nil {
		t.Fatal(err)
	}

	valid, err := testClient.ValidateVaultPassword(testNewVault.ID, testNewVaultPassword)
	if err != nil {
		t.Fatal(err)
	}

	if !valid {
		t.Fatal("vault password validation failed, expected ", testNewVaultPassword)
	}

	vault, err := testClient.GetVault(testNewVault.ID)
	if err != nil {
		t.Fatal(err)
	}

	vault.CreationDate = testNewVault.CreationDate
	vault.ModifiedDate = testNewVault.ModifiedDate

	if !reflect.DeepEqual(testNewVault, vault) {
		t.Fatalf("fetched vault did not match test vault. Expected %#v, got %#v", testNewVault, vault)
	}
}

func test_DeleteVault(t *testing.T) {
	err := testClient.DeleteVault(testNewVault.ID)
	if err != nil {
		t.Fatal(err)
	}
}
