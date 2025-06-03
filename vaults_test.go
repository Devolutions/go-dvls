package dvls

import (
	"reflect"
	"testing"
)

const testNewVaultId string = "eabd3646-acf8-44a4-9ba0-991df147c209"

var testNewVaultPassword string = "5w:mr6kPj"

var testVault Vault = Vault{
	Name:        "go-dvls",
	Description: "Test Vault",
}

var testNewVault Vault = Vault{
	Id:          testNewVaultId,
	Name:        "go-dvls new",
	Description: "Test",
}

func Test_Vaults(t *testing.T) {
	testVault.Id = testVaultId
	t.Run("GetVault", test_GetVault)
	t.Run("NewVault", test_NewVault)
	t.Run("UpdateVault", test_UpdateVault)
	t.Run("DeleteVault", test_DeleteVault)
}

func test_GetVault(t *testing.T) {
	vault, err := testClient.Vaults.Get(testVaultId)
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
	err := testClient.Vaults.New(testNewVault, nil)
	if err != nil {
		t.Fatal(err)
	}

	vault, err := testClient.Vaults.Get(testNewVault.Id)
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

	err := testClient.Vaults.Update(testNewVault, &options)
	if err != nil {
		t.Fatal(err)
	}

	valid, err := testClient.Vaults.ValidatePassword(testNewVault.Id, testNewVaultPassword)
	if err != nil {
		t.Fatal(err)
	}

	if !valid {
		t.Fatal("vault password validation failed, expected ", testNewVaultPassword)
	}

	vault, err := testClient.Vaults.Get(testNewVault.Id)
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
	err := testClient.Vaults.Delete(testNewVault.Id)
	if err != nil {
		t.Fatal(err)
	}
}
