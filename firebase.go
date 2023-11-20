package firebase

import (
	"context"
	"fmt"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

// FirebaseClient is a wrapper around the firebase.App client.
type FirebaseClient struct {
	client *firebase.App
	ctx    context.Context
}

// handleError returns true if the given error is not nil, otherwise false.
func handleError(err error) bool {
	if err != nil {
		return true
	}
	return false
}

// NewFirebaseClient creates a new FirebaseClient based on the given projectID and secretsJSON.
func NewFirebaseClient(ctx context.Context, projectID string, secretsJSON []byte) (*FirebaseClient, error) {
	conf := &firebase.Config{
		ProjectID: projectID,
	}

	opt := option.WithCredentialsJSON(secretsJSON)

	app, err := firebase.NewApp(ctx, conf, opt)
	if handleError(err) {
		return nil, fmt.Errorf("error initializing app: %v", err)
	}

	return &FirebaseClient{client: app, ctx: ctx}, nil
}

// GetDocument returns the document for the given collection and document.
func (f *FirebaseClient) getDocument(collection string, document string) (map[string]interface{}, error) {
	firestoreClient, err := f.client.Firestore(f.ctx)
	if handleError(err) {
		return nil, err
	}
	defer firestoreClient.Close()

	docRef := firestoreClient.Collection(collection).Doc(document)
	doc, err := docRef.Get(f.ctx)
	if handleError(err) {
		return nil, err
	}

	return doc.Data(), nil
}

// DeleteDocument deletes the document for the given collection and document.
func (f *FirebaseClient) deleteDocument(collection string, document string) error {
	firestoreClient, err := f.client.Firestore(f.ctx)
	if handleError(err) {
		return err
	}
	defer firestoreClient.Close()

	docRef := firestoreClient.Collection(collection).Doc(document)
	_, err = docRef.Delete(f.ctx)
	if handleError(err) {
		return err
	}

	return nil
}

// UpsertDocument upserts the document using the given data.
func (f *FirebaseClient) upsertDocument(collection string, document string, data map[string]interface{}) error {
	firestoreClient, err := f.client.Firestore(f.ctx)
	if handleError(err) {
		return err
	}
	defer firestoreClient.Close()

	docRef := firestoreClient.Collection(collection).Doc(document)

	_, err = docRef.Set(f.ctx, data, firestore.MergeAll)
	if handleError(err) {
		return err
	}

	return nil
}
