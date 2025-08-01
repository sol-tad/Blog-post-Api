package repository

import (
	"context"

	"github.com/sol-tad/Blog-post-Api/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type blogRepository struct {
	collection *mongo.Collection
}

func NewBlogRepository(coll *mongo.Collection) domain.BlogRepository {
	return &blogRepository{
		collection: coll,
	}
}

func (r *blogRepository) Create(blog *domain.Blog) error {
	_, err := r.collection.InsertOne(context.Background(), blog)
	return err
}

func (r *blogRepository) GetByID(id string) (*domain.Blog, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var blog domain.Blog
	err = r.collection.FindOne(context.Background(), bson.M{"_id": objID}).Decode(&blog)
	return &blog, err
}

func (r *blogRepository) Update(blog *domain.Blog) error {
	_, err := r.collection.ReplaceOne(
		context.Background(),
		bson.M{"_id": blog.ID},
		blog,
	)
	return err
}

func (r *blogRepository) Delete(id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = r.collection.DeleteOne(context.Background(), bson.M{"_id": objID})
	return err
}

func (r *blogRepository) List(page, limit int, filter domain.BlogFilter) ([]*domain.Blog, int64, error) {
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
	total, err := r.collection.CountDocuments(context.Background(), query)
	if err != nil {
		return nil, 0, err
	}
	
	// Find documents
	cursor, err := r.collection.Find(context.Background(), query, opts)
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

func (r *blogRepository) IncrementViewCount(blogID string) error {
	return r.updateCounter(blogID, "stats.views", 1)
}

func (r *blogRepository) IncrementLikeCount(blogID string) error {
	return r.updateCounter(blogID, "stats.likes", 1)
}

func (r *blogRepository) IncrementDislikeCount(blogID string) error {
	return r.updateCounter(blogID, "stats.dislikes", 1)
}

func (r *blogRepository) DecrementLikeCount(blogID string) error {
	return r.updateCounter(blogID, "stats.likes", -1)
}

func (r *blogRepository) DecrementDislikeCount(blogID string) error {
	return r.updateCounter(blogID, "stats.dislikes", -1)
}

func (r *blogRepository) IncrementCommentCount(blogID string) error {
	return r.updateCounter(blogID, "stats.comments", 1)
}

func (r *blogRepository) DecrementCommentCount(blogID string) error {
	return r.updateCounter(blogID, "stats.comments", -1)
}

// Helper function for atomic counter updates
func (r *blogRepository) updateCounter(blogID, field string, value int) error {
	objID, err := primitive.ObjectIDFromHex(blogID)
	if err != nil {
		return err
	}
	
	_, err = r.collection.UpdateByID(
		context.Background(),
		objID,
		bson.M{"$inc": bson.M{field: value}},
	)
	return err
}