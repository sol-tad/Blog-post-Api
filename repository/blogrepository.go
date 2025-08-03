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
func (b *BlogRepo) GetByAuthor(author string, skip, limit int) ([]*domain.Blog, error) {
    ctx := context.Background()
    filter := bson.M{"author_name": primitive.Regex{Pattern: author, Options: "i"}}
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

func (b *BlogRepo) List(page, limit int, filter domain.BlogFilter) ([]*domain.Blog, int64, error) {
    // Create basic query
    query := bson.M{}
    
    // Add search filter if provided
    if filter.Search != "" {
        query["$or"] = []bson.M{
            {"title": bson.M{"$regex": filter.Search, "$options": "i"}},
            {"author_name": bson.M{"$regex": filter.Search, "$options": "i"}},
        }
    }
    
    // Add tags filter if provided
    if len(filter.Tag) > 0 {
        query["tags"] = bson.M{"$all": filter.Tag}
    }
    
    // Set default sort options
    sortField := "created_at"
    sortOrder := -1  // descending by default
    
    // Apply custom sort if requested
    if filter.SortBy == "popularity" {
        sortField = "stats.views"
    } else if filter.SortBy != "" {
        sortField = filter.SortBy
    }
    
    // Create options
    opts := options.Find().
        SetSort(bson.D{{Key: sortField, Value: sortOrder}}).
        SetSkip(int64((page - 1) * limit)).
        SetLimit(int64(limit))
    
    // Get total count
    total, err := b.collection.CountDocuments(context.TODO(), query)
    if err != nil {
        return nil, 0, err
    }
    
    // Find documents
    cursor, err := b.collection.Find(context.TODO(), query, opts)
    if err != nil {
        return nil, 0, err
    }
    
    // Decode results
    var blogs []*domain.Blog
    if err = cursor.All(context.TODO(), &blogs); err != nil {
        return nil, 0, err
    }
    
    return blogs, total, nil
}