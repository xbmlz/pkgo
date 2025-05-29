package gormx

import (
	"context"

	"gorm.io/gorm"
)

type Service[T any] struct {
	db *gorm.DB
}

func NewService[T any](db *gorm.DB) *Service[T] {
	return &Service[T]{db: db}
}

func (s *Service[T]) Count(ctx context.Context, conditions ...any) (int64, error) {
	var (
		count  int64
		entity T
	)
	err := s.db.WithContext(ctx).Model(entity).Where(conditions).Count(&count).Error
	return count, err
}

func (s *Service[T]) Create(ctx context.Context, entity T) error {
	return s.db.WithContext(ctx).Create(&entity).Error
}

func (s *Service[T]) CreateOrUpdate(ctx context.Context, entity T, conditions ...any) error {
	if err := s.db.WithContext(ctx).Where(conditions).First(&entity).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return s.Create(ctx, entity)
		}
		return err
	}
	return s.Update(ctx, entity, conditions)
}

func (s *Service[T]) FindByID(ctx context.Context, id uint) (*T, error) {
	var entity T
	err := s.db.WithContext(ctx).First(&entity, id).Error
	return &entity, err
}

func (s *Service[T]) FindOne(ctx context.Context, conditions ...any) (*T, error) {
	var entity T
	err := s.db.WithContext(ctx).First(&entity, conditions).Error
	return &entity, err
}

func (s *Service[T]) Find(ctx context.Context, conditions ...any) ([]T, error) {
	var entities []T
	err := s.db.WithContext(ctx).Find(&entities, conditions).Error
	return entities, err
}

func (s *Service[T]) Update(ctx context.Context, entity T, conditions ...any) error {
	return s.db.WithContext(ctx).Model(&entity).Where(conditions).Updates(entity).Error
}

func (s *Service[T]) UpdateByID(ctx context.Context, id uint, entity T) error {
	return s.db.WithContext(ctx).Model(&entity).Where("id = ?", id).Updates(entity).Error
}

func (s *Service[T]) Delete(ctx context.Context, entity T, conditions ...any) error {
	return s.db.WithContext(ctx).Delete(&entity, conditions).Error
}

func (s *Service[T]) DeleteByID(ctx context.Context, id uint) error {
	var entity T
	return s.db.WithContext(ctx).Delete(&entity, id).Error
}

func (s *Service[T]) Page(ctx context.Context, page PageQuery, conditions ...any) (PageResult[T], error) {
	var entity T
	m := s.db.WithContext(ctx).Model(entity)
	if len(conditions) > 0 {
		m = m.Where(conditions)
	}
	if len(page.Orders) > 0 {
		for _, order := range page.Orders {
			if order.Desc {
				m = m.Order(order.Column + " desc")
			} else {
				m = m.Order(order.Column + " asc")
			}
		}
	}
	var total int64
	err := m.Count(&total).Error
	if err != nil {
		return PageResult[T]{}, err
	}
	var records []T
	err = m.Offset((page.Page - 1) * page.PageSize).Limit(page.PageSize).Find(&records).Error
	if err != nil {
		return PageResult[T]{}, err
	}
	return PageResult[T]{
		Page:     page.Page,
		PageSize: page.PageSize,
		Total:    int(total),
		Records:  records,
	}, nil
}
