package repository

import (
	"axolotl-cloud/infra/cache"
	"axolotl-cloud/infra/settings"
	"axolotl-cloud/internal/app/model"
	"sync"
	"time"

	"gorm.io/gorm"
)

var defaultSettings = []model.Setting{
	{Key: settings.JobTimeout, Value: "1800"},
	{Key: settings.Language, Value: "en"},
}

type SettingRepository struct {
	DB    *gorm.DB
	cache *cache.TypedCache[map[model.SettingKey]model.Setting]
	mu    sync.Mutex
}

func NewSettingRepository(db *gorm.DB) *SettingRepository {
	r := &SettingRepository{DB: db}
	r.initDefaults()
	r.initCache()
	return r
}

func (r *SettingRepository) initDefaults() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	var existing []model.Setting
	if err := r.DB.Find(&existing).Error; err != nil {
		return err
	}

	existingMap := make(map[model.SettingKey]struct{})
	for _, s := range existing {
		existingMap[s.Key] = struct{}{}
	}

	for _, setting := range defaultSettings {
		if _, exists := existingMap[setting.Key]; !exists {
			if err := r.DB.Create(&setting).Error; err != nil {
				return err
			}
		}
	}
	return nil
}

func (r *SettingRepository) initCache() {
	r.cache = cache.NewTypedCache(
		make(map[model.SettingKey]model.Setting),
		12*time.Hour,
		func(...any) (map[model.SettingKey]model.Setting, error) {
			return r.loadAllSettings()
		},
	)
}

func (r *SettingRepository) loadAllSettings() (map[model.SettingKey]model.Setting, error) {
	var settings []model.Setting
	if err := r.DB.Find(&settings).Error; err != nil {
		return nil, err
	}

	result := make(map[model.SettingKey]model.Setting, len(settings))
	for _, s := range settings {
		result[s.Key] = s
	}
	return result, nil
}

func (r *SettingRepository) Save(setting *model.Setting) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if err := r.DB.Save(setting).Error; err != nil {
		return err
	}

	r.cache.Update(func(curr map[model.SettingKey]model.Setting) map[model.SettingKey]model.Setting {
		newMap := make(map[model.SettingKey]model.Setting, len(curr))
		for k, v := range curr {
			newMap[k] = v
		}
		newMap[setting.Key] = *setting
		return newMap
	})
	return nil
}

func (r *SettingRepository) GetAll() ([]model.Setting, error) {
	data, err := r.cache.Get()
	if err != nil {
		return nil, err
	}

	settings := make([]model.Setting, 0, len(data))
	for _, s := range data {
		settings = append(settings, s)
	}
	return settings, nil
}

func (r *SettingRepository) GetByKey(key model.SettingKey) (*model.Setting, error) {
	data, err := r.cache.Get()
	if err != nil {
		return nil, err
	}

	if setting, exists := data[key]; exists {
		return &setting, nil
	}
	return nil, gorm.ErrRecordNotFound
}

func (r *SettingRepository) RemoveByKey(key string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if err := r.DB.Delete(&model.Setting{}, "key = ?", key).Error; err != nil {
		return err
	}

	r.cache.Update(func(curr map[model.SettingKey]model.Setting) map[model.SettingKey]model.Setting {
		newMap := make(map[model.SettingKey]model.Setting, len(curr))
		for k, v := range curr {
			if k != model.SettingKey(key) {
				newMap[k] = v
			}
		}
		return newMap
	})
	return nil
}
