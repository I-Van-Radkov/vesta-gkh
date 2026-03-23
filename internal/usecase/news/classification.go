package news

import (
	"fmt"
	"math"
	"sort"
	"strings"
	"time"
)

func (p *NewsParser) classify(title, desc, content string) (bool, float64, string) {
	p.loadCache()

	text := normalize(title + " " + desc + " " + content)
	textLower := strings.ToLower(text)
	titleLower := strings.ToLower(title)

	score := 0.0

	// Ключевые слова
	for _, kw := range p.keywords {
		if !kw.IsActive {
			continue
		}
		kwLower := strings.ToLower(kw.Keyword)

		countTitle := float64(strings.Count(titleLower, kwLower))
		countDesc := float64(strings.Count(strings.ToLower(desc), kwLower))
		countCont := float64(strings.Count(strings.ToLower(content), kwLower))

		if countTitle+countDesc+countCont == 0 {
			continue
		}

		raw := countTitle*p.cfg.TitleWeight + countDesc*p.cfg.DescWeight + countCont*p.cfg.ContentWeight
		contrib := kw.Weight * math.Min(raw, p.cfg.MaxContribKeyword)
		if raw > 4 {
			contrib *= 1 / (1 + p.cfg.DampeningFactor*math.Log1p(raw-4))
		}

		if contrib > 1 {
			score += contrib
		}
	}

	// Позитивные regex
	for _, r := range p.regexRules {
		if !r.IsActive || r.Compiled == nil {
			continue
		}
		if r.Compiled.MatchString(textLower) {
			contrib := math.Min(r.BonusScore, p.cfg.MaxContribRegex)
			score += contrib
		}
	}

	// Категории
	catScores := make([]struct {
		name  string
		score float64
	}, 0, len(p.categories))
	for _, cat := range p.categories {
		hits := 0
		for _, word := range cat.Keywords {
			if strings.Contains(textLower, strings.ToLower(word)) {
				hits++
			}
		}
		if hits >= cat.MinHitsForBonus {
			contrib := math.Min(cat.BonusPerHit*float64(hits), p.cfg.MaxContribCategory)
			score += contrib
			catScores = append(catScores, struct {
				name  string
				score float64
			}{cat.Name, contrib})
		}
	}

	// Big bonus
	if len(catScores) >= 2 {
		sort.Slice(catScores, func(i, j int) bool { return catScores[i].score > catScores[j].score })
		big := catScores[0].score * 0.75
		score += big
	}

	// Негативные regex
	for _, n := range p.negRules {
		if !n.IsActive || n.Compiled == nil {
			continue
		}
		if n.Compiled.MatchString(textLower) {
			score += n.Penalty
		}
	}

	// Обязательный фильтр по заголовку
	titleOnlyScore := 0.0
	for _, kw := range p.keywords {
		if strings.Contains(titleLower, strings.ToLower(kw.Keyword)) {
			titleOnlyScore += kw.Weight * p.cfg.TitleWeight
		}
	}

	ok := score >= p.cfg.Threshold
	debug := fmt.Sprintf("total=%.1f | titleOnly=%.1f", score, titleOnlyScore)

	return ok, score, debug
}

func (p *NewsParser) loadCache() {
	p.mu.Lock()
	defer p.mu.Unlock()

	if time.Since(p.lastCache) < 5*time.Minute && len(p.keywords) > 0 {
		return
	}

	p.keywords, _ = p.newsRepo.GetKeywordsList()
	p.regexRules, _ = p.newsRepo.GetRegexRulesList()
	p.negRules, _ = p.newsRepo.GetNegativeRegexList()
	p.categories, _ = p.newsRepo.GetCategoryWithKeywords()
	p.lastCache = time.Now()
}
