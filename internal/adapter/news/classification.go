package news

import (
	"context"
	"fmt"
	"log"
	"regexp"

	"github.com/I-Van-Radkov/vesta-gkh/internal/models"
)

func (r *NewsRepo) GetKeywordsList() ([]*models.GKHKeyword, error) {
	rows, err := r.db.Query(context.Background(),
		`SELECT id, keyword, weight, category, is_active, created_at, updated_at
		 FROM gkh_keywords
		 WHERE is_active = true`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*models.GKHKeyword
	for rows.Next() {
		var kw models.GKHKeyword
		err := rows.Scan(
			&kw.ID, &kw.Keyword, &kw.Weight, &kw.Category,
			&kw.IsActive, &kw.CreatedAt, &kw.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		list = append(list, &kw)
	}
	return list, rows.Err()
}

func (r *NewsRepo) GetRegexRulesList() ([]*models.GKHRegexRule, error) {
	rows, err := r.db.Query(context.Background(),
		`SELECT id, pattern, bonus_score, is_active, created_at, updated_at
		 FROM gkh_regex_rules
		 WHERE is_active = true`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rules []*models.GKHRegexRule
	for rows.Next() {
		var rule models.GKHRegexRule
		var pattern string
		err := rows.Scan(
			&rule.ID, &pattern, &rule.BonusScore,
			&rule.IsActive, &rule.CreatedAt, &rule.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		rule.PatternStr = pattern
		rule.Compiled, err = regexp.Compile(pattern)
		if err != nil {
			log.Printf("bad regex pattern %q → skipped: %v", pattern, err)
			continue
		}
		rules = append(rules, &rule)
	}
	return rules, rows.Err()
}

func (r *NewsRepo) GetNegativeRegexList() ([]*models.GKHNegativeRegex, error) {
	rows, err := r.db.Query(context.Background(),
		`SELECT id, pattern, penalty, description, is_active, created_at, updated_at
		 FROM gkh_negative_regex
		 WHERE is_active = true`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*models.GKHNegativeRegex
	for rows.Next() {
		var nr models.GKHNegativeRegex
		var pattern string
		err := rows.Scan(
			&nr.ID, &pattern, &nr.Penalty, &nr.Description,
			&nr.IsActive, &nr.CreatedAt, &nr.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		nr.PatternStr = pattern
		nr.Compiled, err = regexp.Compile(pattern)
		if err != nil {
			log.Printf("bad negative regex %q → skipped: %v", pattern, err)
			continue
		}
		list = append(list, &nr)
	}
	return list, rows.Err()
}

func (r *NewsRepo) GetCategoryWithKeywords() ([]*models.GHKCategoryWithKeywords, error) {
	rows, err := r.db.Query(context.Background(), `
		SELECT 
			c.id, c.name, c.bonus_per_hit, c.min_hits_for_bonus, c.big_bonus, c.description,
			COALESCE(ARRAY_AGG(k.keyword) FILTER (WHERE k.keyword IS NOT NULL), '{}') AS keywords
		FROM gkh_categories c
		LEFT JOIN gkh_category_keywords ck ON ck.category_id = c.id
		LEFT JOIN gkh_keywords k ON ck.keyword_id = k.id AND k.is_active
		GROUP BY c.id
		HAVING COUNT(k.id) > 0
	`)
	if err != nil {
		return nil, fmt.Errorf("query categories with keywords: %w", err)
	}
	defer rows.Close()

	var cats []*models.GHKCategoryWithKeywords
	for rows.Next() {
		var cat models.GHKCategoryWithKeywords
		var keywordsArray []string
		err := rows.Scan(
			&cat.ID, &cat.Name, &cat.BonusPerHit, &cat.MinHitsForBonus,
			&cat.BigBonus, &cat.Description, &keywordsArray,
		)
		if err != nil {
			return nil, fmt.Errorf("scan category: %w", err)
		}
		cat.Keywords = keywordsArray
		cats = append(cats, &cat)
	}
	return cats, rows.Err()
}
