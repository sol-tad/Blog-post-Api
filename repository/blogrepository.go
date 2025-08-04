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

type BlogRepo struct{
	collection *mongo.Collection
	context context.Context
}

func NewBlogRepo(coll *mongo.Collection) usecase.IBlogRepo {
	ctx := context.Background()
	return &BlogRepo{
		collection: coll,
		context: ctx,
	}
}

func (b *BlogRepo) StoreBlog (blog *domain.Blog)error{
	_, err := b.collection.InsertOne(b.context, blog)
	return err
}

func (b *BlogRepo) RetriveAll ()[]domain.Blog{
	var results []domain.Blog
	filter := bson.D{}
	cursor , err := b.collection.Find(b.context, filter)
		if err != nil{
		fmt.Println("Error finding blogs:", err)
		return []domain.Blog{} 
	}

	err = cursor.All(b.context, &results)
		if err != nil{
		fmt.Println("Error decoding Blogs:", err)
		return []domain.Blog{} 
	}
	return results
}
func (b *BlogRepo) ViewBlogByID(blogID primitive.ObjectID) *domain.Blog{
	var result domain.Blog
	filter := bson.M{"_id":blogID}
	err := b.collection.FindOne(b.context, filter).Decode(&result)
	if err != nil{
		fmt.Println("blog not found or error decoding:", err)
		return nil
	}
	return &result
}
func (b *BlogRepo) UpdateBlog(id primitive.ObjectID, updatedBlog *domain.Blog) error{
	filter := bson.M{"_id":id}
	updated := bson.M{
		"$set":bson.M{
	"content":updatedBlog.Content,
	"id" : updatedBlog.ID,
	"title" : updatedBlog.Title,
	"tags":updatedBlog.Tags,
	"created_at ":updatedBlog.CreatedAt,
	"updated_at" : updatedBlog.UpdatedAt,
	"author_id" : updatedBlog.AuthorID,
	"author_name": updatedBlog.AuthorName,
	"stats" : updatedBlog.Stats,
			
		},
	}

	_ , err := b.collection.UpdateOne(b.context,filter,updated)
	return err
}

func (b *BlogRepo) DeleteBlog(id primitive.ObjectID) error{

	filter := bson.M{"_id": id}
	_, err := b.collection.DeleteOne(b.context, filter)
	return err
}


// on this part, i have just make naming modifications inorder to be consistent. so , you can just focus on the implementation.


func (b *BlogRepo) List(page, limit int, filter domain.BlogFilter) ([]*domain.Blog, int64, error) {
	// Build filter query
	query := bson.M{}
	
	if filter.Search != "" {
		query["$or"] = bson.A{
			bson.M{"title": primitive.Regex{Pattern: filter.Search, Options: "i"}},
			bson.M{"author_name": primitive.Regex{Pattern: filter.Search, Options: "i"}},
		}
	}
	
	if filter.Author != "" {
		query["author_name"] = primitive.Regex{Pattern: filter.Author, Options: "i"}
	}
	
	if len(filter.Tags) > 0 {
		query["tags"] = bson.M{"$all": filter.Tags}
	}
	
	if !filter.StartDate.IsZero() && !filter.EndDate.IsZero() {
		query["created_at"] = bson.M{
			"$gte": filter.StartDate,
			"$lte": filter.EndDate,
		}
	}
	
	// Sorting
	sortOption := -1 // default descending
	if filter.SortOrder == "asc" {
		sortOption = 1
	}
	
	sort := bson.D{{Key: filter.SortBy, Value: sortOption}}
	
	// Pagination options
	opts := options.Find().
		SetSort(sort).
		SetSkip(int64((page - 1) * limit)).
		SetLimit(int64(limit))
	
	// Get total count
	total, err := b.collection.CountDocuments(context.Background(), query)
	if err != nil {
		return nil, 0, err
	}
	
	// Find documents
	cursor, err := b.collection.Find(context.Background(), query, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(context.Background())
	
	var blogs []*domain.Blog
	if err = cursor.All(context.Background(), &blogs); err != nil {
		return nil, 0, err
	}
	
	return blogs, total, nil
}

func (b *BlogRepo) IncrementViewCount(blogID string) error {
	return b.updateCounter(blogID, "stats.views", 1)
}

func (b *BlogRepo) IncrementLikeCount(blogID string) error {
	return b.updateCounter(blogID, "stats.likes", 1)
}

func (b *BlogRepo) IncrementDislikeCount(blogID string) error {
	return b.updateCounter(blogID, "stats.dislikes", 1)
}

func (b *BlogRepo) DecrementLikeCount(blogID string) error {
	return b.updateCounter(blogID, "stats.likes", -1)
}

func (b *BlogRepo) DecrementDislikeCount(blogID string) error {
	return b.updateCounter(blogID, "stats.dislikes", -1)
}

func (b *BlogRepo) IncrementCommentCount(blogID string) error {
	return b.updateCounter(blogID, "stats.comments", 1)
}

func (b *BlogRepo) DecrementCommentCount(blogID string) error {
	return b.updateCounter(blogID, "stats.comments", -1)
}

// Helper function for atomic counter updates
func (b *BlogRepo) updateCounter(blogID, field string, value int) error {
	objID, err := primitive.ObjectIDFromHex(blogID)
	if err != nil {
		return err
	}
	
	_, err = b.collection.UpdateByID(
		context.Background(),
		objID,
		bson.M{"$inc": bson.M{field: value}},
	)
	return err
}