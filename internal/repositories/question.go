package repositories

import (
	"github.com/sadia-54/qstack-backend/internal/models/domains"
	"gorm.io/gorm"
)

type QuestionRepository struct {
	db *gorm.DB
}

func NewQuestionRepository(db *gorm.DB) *QuestionRepository {
	return &QuestionRepository{db: db}
}

func (r *QuestionRepository) Create(question *domains.Question) error {
	return r.db.Create(question).Error
}

func (r *QuestionRepository) Update(question *domains.Question) error {
	return r.db.Session(&gorm.Session{FullSaveAssociations: true}).Updates(question).Error
}

func (r *QuestionRepository) FindByID(id int64) (*domains.Question, error) {
	var question domains.Question
	err := r.db.Preload("User").
		Preload("Tags").
		First(&question, id).Error

	if err != nil {
		return nil, err
	}

	return &question, nil
}

func (r *QuestionRepository) GetFeed(
	search string,
	tag string,
	sort string,
	limit int,
	offset int,
) ([]domains.Question, error) {

	var questions []domains.Question

	query := r.db.
		Model(&domains.Question{}).
		Preload("User").
		Preload("Tags")

	// Search by title
	if search != "" {
		query = query.Where("title ILIKE ?", "%"+search+"%")
	}

	// Filter by tag
	if tag != "" {
		query = query.
			Joins("JOIN question_tags qt ON qt.question_id = questions.id").
			Joins("JOIN tags t ON t.id = qt.tag_id").
			Where("t.name = ?", tag)
	}

	// Sorting
	switch sort {
	case "votes":
		query = query.Order("vote_count DESC")
	case "date":
		query = query.Order("created_at DESC")
	default:
		query = query.Order("created_at DESC")
	}

	err := query.
		Limit(limit).
		Offset(offset).
		Find(&questions).Error

	return questions, err
}

func (r *QuestionRepository) Delete(id int64) error {
	return r.db.Delete(&domains.Question{}, id).Error
}

func (r *QuestionRepository) UpdateVoteCount(questionID int64, diff int) error {

	return r.db.
		Model(&domains.Question{}).
		Where("id = ?", questionID).
		UpdateColumn("vote_count", gorm.Expr("vote_count + ?", diff)).
		Error
}

func (r *QuestionRepository) GetMyFeed(
	userID int64,
	limit int,
	offset int,
) ([]domains.Question, error) {

	var questions []domains.Question

	subQuery := `
	    SELECT qt.tag_id
	    FROM questions q2
	    JOIN question_tags qt ON qt.question_id = q2.id
	    WHERE q2.user_id = ?

	    UNION

	    SELECT qt.tag_id
	    FROM question_votes v
	    JOIN question_tags qt ON qt.question_id = v.question_id
	    WHERE v.user_id = ?

	    UNION

	    SELECT qt.tag_id
	    FROM answers a
	    JOIN question_tags qt ON qt.question_id = a.question_id
	    WHERE a.user_id = ?
	`

	err := r.db.
		Model(&domains.Question{}).
		Distinct("questions.*").
		Preload("User").
		Preload("Tags").
		Joins("JOIN question_tags qt ON qt.question_id = questions.id").
		Where("qt.tag_id IN ("+subQuery+")", userID, userID, userID).
		Where("questions.user_id != ?", userID).
		Order("questions.created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&questions).Error

	return questions, err
}