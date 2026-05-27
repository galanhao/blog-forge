package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/galanhao/blog-forge/internal/config"
	"github.com/galanhao/blog-forge/internal/site"
	"github.com/galanhao/blog-forge/internal/util"
)

const (
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
	fmt.Println("blog-forge - static site generator")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  blog-forge build [-site DIR] [-config FILE] [-themes DIR] [-dist DIR]")
	fmt.Println("  blog-forge new <title>")
	fmt.Println("  blog-forge init")
}

func runBuild() error {
	buildCmd := flag.NewFlagSet("build", flag.ExitOnError)
	siteDir := buildCmd.String("site", ".", "site root directory (contains content/, assets/, static/)")
	configFile := buildCmd.String("config", "_config.yml", "path to config file")
	themesDir := buildCmd.String("themes", defaultThemes, "path to themes directory")
	distDir := buildCmd.String("dist", "dist", "output directory")

	buildCmd.Parse(os.Args[2:])

	cfg, err := config.Load(*configFile)
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	builder, err := site.New(cfg, *siteDir, *themesDir)
	if err != nil {
		return fmt.Errorf("init builder: %w", err)
	}

	if err := builder.Build(*distDir); err != nil {
		return fmt.Errorf("build: %w", err)
	}

	fmt.Printf("✅ Site built to %s/\n", *distDir)
	return nil
}

func runNew() error {
	if len(os.Args) < 3 {
		return fmt.Errorf("usage: blog-forge new <title>")
	}
	title := os.Args[2]
	return createPost(title)
}

func createPost(title string) error {
	slug := util.Slugify(title)
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

func nowDate() string {
	now := time.Now()
	return fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d",
		now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second())
}

func runInit() error {
	return fmt.Errorf("init command not yet implemented")
}
