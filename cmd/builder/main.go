package main

import (
	"fmt"
	"os"

	"github.com/galanhao/blog-builder/internal/config"
	"github.com/galanhao/blog-builder/internal/site"
)

const (
	defaultConfig = "_config.yml"
	defaultDist   = "dist"
	defaultThemes = "themes"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	cmd := os.Args[1]
	switch cmd {
	case "build":
		if err := runBuild(); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
	case "new":
		if err := runNew(); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
	case "init":
		if err := runInit(); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n", cmd)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("blog-builder - static site generator")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  blog-builder build    Build the site")
	fmt.Println("  blog-builder new      Create a new post")
	fmt.Println("  blog-builder init     Initialize a new site")
}

func runBuild() error {
	cfg, err := config.Load(defaultConfig)
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	builder, err := site.New(cfg, defaultThemes)
	if err != nil {
		return fmt.Errorf("init builder: %w", err)
	}

	if err := builder.Build(defaultDist); err != nil {
		return fmt.Errorf("build: %w", err)
	}

	fmt.Printf("✅ Site built to %s/\n", defaultDist)
	return nil
}

func runNew() error {
	if len(os.Args) < 3 {
		return fmt.Errorf("usage: blog-builder new <title>")
	}
	title := os.Args[2]
	return createPost(title)
}

func createPost(title string) error {
	slug := slugify(title)
	filename := fmt.Sprintf("content/posts/%s.md", slug)

	tmpl := fmt.Sprintf(`---
title: "%s"
date: %s
tags: []
categories: []
excerpt: ""
cover_image: ""
---

`, title, nowDate())

	if err := os.WriteFile(filename, []byte(tmpl), 0o644); err != nil {
		return fmt.Errorf("create post: %w", err)
	}

	fmt.Printf("📝 Created %s\n", filename)
	return nil
}

func slugify(s string) string {
	result := make([]rune, 0, len(s))
	for _, r := range s {
		if r >= 'a' && r <= 'z' || r >= '0' && r <= '9' {
			result = append(result, r)
		} else if r >= 'A' && r <= 'Z' {
			result = append(result, r+32) // to lower
		} else if r == ' ' || r == '_' {
			result = append(result, '-')
		} else if r >= 0x4e00 { // CJK characters - keep
			result = append(result, r)
		}
	}
	return string(result)
}

func nowDate() string {
	return fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d",
		2026, 5, 25, 0, 0, 0) // TODO: use time.Now()
}

func runInit() error {
	return fmt.Errorf("init command not yet implemented")
}
