package model

import "time"

func DefaultSiteSetting(now time.Time) SiteSetting {
	shortName := "字里行间"
	titleSuffix := "字里行间"
	technologyText := "Built with Vue 3 and Go"
	startYear := now.UTC().Year()
	return SiteSetting{
		ID:              1,
		SiteName:        "字里行间",
		SiteShortName:   &shortName,
		SiteDescription: "记录开发、阅读与生活中的思考。",
		TitleSuffix:     &titleSuffix,
		CopyrightName:   "字里行间",
		StartYear:       &startYear,
		ShowTechnology:  true,
		TechnologyText:  &technologyText,
	}
}

func DefaultHomepageConfig() HomepageConfig {
	return HomepageConfig{Modules: []HomepageModule{
		{
			Type: HomepageModuleHero, Enabled: true, SortOrder: 10,
			Hero: &HeroModuleConfig{
				Eyebrow: "ZERVAN'S JOURNAL", Title: "在字里行间，保存思想的回声。",
				HighlightText: "保存思想的回声。", Description: "记录开发、阅读与生活中的思考。把遇到的问题讲清楚，也把值得留存的经验认真写下来。",
				Layout:          HeroLayoutCenter,
				PrimaryButton:   HomepageButton{Enabled: true, Text: "开始阅读", URL: "/archive", LinkType: LinkTypeInternal},
				SecondaryButton: HomepageButton{Enabled: true, Text: "关于本站", URL: "/about", LinkType: LinkTypeInternal},
			},
		},
		{
			Type: HomepageModuleAbout, Enabled: true, SortOrder: 20,
			About: &AboutModuleConfig{
				Title: "关于我", Description: "保持好奇，持续记录。",
				Content:       "这里记录开发、阅读与生活中的思考。\n\n写作让零散的经验变成可以反复回看的文字。",
				ImagePosition: AboutImageNone,
			},
		},
		{
			Type: HomepageModuleFeaturedArticles, Enabled: true, SortOrder: 30,
			FeaturedArticles: &FeaturedArticlesModuleConfig{Title: "推荐文章", Description: "精选内容", Limit: 3},
		},
		{
			Type: HomepageModuleLatestArticles, Enabled: true, SortOrder: 40,
			LatestArticles: &LatestArticlesModuleConfig{
				Title: "最新文章", Limit: 10, ShowSummary: true, ShowDate: true, ShowCommentCount: true, ShowViewAll: true,
			},
		},
		{
			Type: HomepageModuleTechStack, Enabled: false, SortOrder: 50,
			TechStack: &TechStackModuleConfig{
				Title: "技术栈", Description: "构建这个博客使用的工具。",
				Items: []TechItem{
					{Name: "Vue 3", IsVisible: true, SortOrder: 10},
					{Name: "Go", IsVisible: true, SortOrder: 20},
					{Name: "PostgreSQL", IsVisible: true, SortOrder: 30},
				},
			},
		},
		{
			Type: HomepageModuleSocialLinks, Enabled: true, SortOrder: 60,
			SocialLinks: &SocialLinksModuleConfig{Title: "找到我"},
		},
	}}
}

func DefaultNavigationItems() []NavigationItem {
	return []NavigationItem{
		{Name: "首页", URL: "/", LinkType: LinkTypeInternal, IsVisible: true, SortOrder: 10},
		{Name: "归档", URL: "/archive", LinkType: LinkTypeInternal, IsVisible: true, SortOrder: 20},
		{Name: "关于", URL: "/about", LinkType: LinkTypeInternal, IsVisible: true, SortOrder: 30},
		{Name: "GitHub", URL: "https://github.com/Zervan29131", LinkType: LinkTypeExternal, OpenInNewTab: true, IsVisible: true, SortOrder: 40},
	}
}

func DefaultSocialLinks() []SocialLink {
	return []SocialLink{
		{Platform: SocialPlatformGitHub, DisplayName: "GitHub", URL: "https://github.com/Zervan29131", IsVisible: true, SortOrder: 10},
	}
}
