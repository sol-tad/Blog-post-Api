package repository

import (
	"context"
	"fmt"

	"github.com/sol-tad/Blog-post-Api/domain"
	"github.com/sol-tad/Blog-post-Api/usecase"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// BlogRepo implements the usecase.IBlogRepo interface
// and provides MongoDB-based storage for blogs.
type BlogRepo struct {
	collection *mongo.Collection // MongoDB collection for blogs
	context    context.Context   // Context for MongoDB operations
}

// NewBlogRepo creates a new BlogRepo instance with the given MongoDB collection.
// It also initializes a background context.
func NewBlogRepo(coll *mongo.Collection) usecase.IBlogRepo {
	ctx := context.Background()
	return &BlogRepo{
		collection: coll,
		context:    ctx,
	}
}

// StoreBlog inserts a new blog document into the collection.
// If successful, it sets the inserted document's ID back to the blog object.
func (b *BlogRepo) StoreBlog(blog *domain.Blog) error {
	result, err := b.collection.InsertOne(b.context, blog)
	if err != nil {
		return err
	}

	// Convert InsertedID to ObjectID and assign to blog.ID
	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		blog.ID = oid
	}

	return nil
}

// RetriveAll returns all blog documents from the collection.
// Returns an empty slice if any error occurs.
func (b *BlogRepo) RetriveAll() []domain.Blog {
	var results []domain.Blog
	filter := bson.D{} // empty filter to match all documents

	cursor, err := b.collection.Find(b.context, filter)
	if err != nil {
		fmt.Println("Error finding blogs:", err)
		return []domain.Blog{}
	}

	err = cursor.All(b.context, &results)
	if err != nil {
		fmt.Println("Error decoding Blogs:", err)
		return []domain.Blog{}
	}

	return results
}

// ViewBlogByID finds and returns a blog by its ObjectID.
// Returns nil if not found or on error.
func (b *BlogRepo) ViewBlogByID(blogID primitive.ObjectID) *domain.Blog {
	var result domain.Blog
	filter := bson.M{"_id": blogID}

	err := b.collection.FindOne(b.context, filter).Decode(&result)
	if err != nil {
		fmt.Println("blog not found or error decoding:", err)
		return nil
	}

	return &result
}

// GetByAuthor retrieves blogs by author name (case-insensitive search),
// with pagination support using skip and limit.
func (b *BlogRepo) GetByAuthor(author string, skip, limit int) ([]*domain.Blog, error) {
	ctx := context.Background()

	// Use regex to do case-insensitive match on author_name field
	filter := bson.M{"author_name": primitive.Regex{Pattern: author, Options: "i"}}

	// Set pagination options
	opts := options.Find().
		SetSkip(int64(skip)).
		SetLimit(int64(limit))

	cursor, err := b.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var blogs []*domain.Blog
	if err = cursor.All(ctx, &blogs); err != nil {
		return nil, err
	}

	return blogs, nil
}

// UpdateBlog updates an existing blog document identified by id with the given updatedBlog fields.
func (b *BlogRepo) UpdateBlog(id primitive.ObjectID, updatedBlog *domain.Blog) error {
	filter := bson.M{"_id": id}

	// Prepare the update document with $set operator
	updated := bson.M{
		"$set": bson.M{
			"content":     updatedBlog.Content,
			"id":          updatedBlog.ID,
			"title":       updatedBlog.Title,
			"tags":        updatedBlog.Tags,
			"created_at ": updatedBlog.CreatedAt,
			"updated_at":  updatedBlog.UpdatedAt,
			"author_id":   updatedBlog.AuthorID,
			"author_name": updatedBlog.AuthorName,
			"stats":       updatedBlog.Stats,
		},
	}

	_, err := b.collection.UpdateOne(b.context, filter, updated)
	return err
}

// DeleteBlog deletes a blog document by its ObjectID.
func (b *BlogRepo) DeleteBlog(id primitive.ObjectID) error {
	filter := bson.M{"_id": id}
	_, err := b.collection.DeleteOne(b.context, filter)
	return err
}

// List returns a paginated list of blogs filtered by search term and tags,
// sorted by a specified field, along with the total count matching the filters.
func (b *BlogRepo) List(page, limit int, filter domain.BlogFilter) ([]*domain.Blog, int64, error) {
	// Initialize an empty query map
	query := bson.M{}

	// Add search filter on title or author_name (case-insensitive) if provided
	if filter.Search != "" {
		query["$or"] = []bson.M{
			{"title": bson.M{"$regex": filter.Search, "$options": "i"}},
			{"author_name": bson.M{"$regex": filter.Search, "$options": "i"}},
		}
	}

	// Add tag filtering if tags are specified
	if len(filter.Tag) > 0 {
		query["tags"] = bson.M{"$all": filter.Tag}
	}

	// Default sorting by created_at descending
	sortField := "created_at"
	sortOrder := -1

	// Override sort field if specified
	if filter.SortBy == "popularity" {
		sortField = "stats.views"
	} else if filter.SortBy != "" {
		sortField = filter.SortBy
	}

	// Set find options: sort, skip and limit for pagination
	opts := options.Find().
		SetSort(bson.D{{Key: sortField, Value: sortOrder}}).
		SetSkip(int64((page - 1) * limit)).
		SetLimit(int64(limit))

	// Count total documents matching the query (for pagination info)
	total, err := b.collection.CountDocuments(context.TODO(), query)
	if err != nil {
		return nil, 0, err
	}

	// Execute find query with options
	cursor, err := b.collection.Find(context.TODO(), query, opts)
	if err != nil {
		return nil, 0, err
	}

	// Decode all results into slice
	var blogs []*domain.Blog
	if err = cursor.All(context.TODO(), &blogs); err != nil {
		return nil, 0, err
	}

	// Return the list and total count
	return blogs, total, nil
}
