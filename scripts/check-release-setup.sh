#!/bin/bash

echo "🔍 Checking GitHub release setup..."

# Check if we're in a git repository
if ! git rev-parse --git-dir > /dev/null 2>&1; then
    echo "❌ Not in a git repository"
    exit 1
fi

# Check remote URL
echo "📡 Remote URL:"
git remote -v

# Check current branch
echo "🌿 Current branch:"
git branch --show-current

# Check if tag exists
TAG=$(git describe --tags --exact-match 2>/dev/null || echo "")
if [ -n "$TAG" ]; then
    echo "🏷️  Current tag: $TAG"
else
    echo "⚠️  No exact tag match found"
fi

# Check if tag is pushed
if [ -n "$TAG" ]; then
    if git ls-remote --tags origin | grep -q "refs/tags/$TAG"; then
        echo "✅ Tag $TAG is pushed to remote"
    else
        echo "❌ Tag $TAG is not pushed to remote"
        echo "   Run: git push origin $TAG"
    fi
fi

# Check workflow files
echo "📋 Workflow files:"
if [ -f ".github/workflows/release.yml" ]; then
    echo "✅ Release workflow exists"
else
    echo "❌ Release workflow missing"
fi

# Check permissions in workflow
if [ -f ".github/workflows/release.yml" ]; then
    if grep -q "permissions:" .github/workflows/release.yml; then
        echo "✅ Permissions configured in workflow"
    else
        echo "⚠️  No explicit permissions in workflow"
    fi
fi

echo ""
echo "🔧 Common fixes:"
echo "1. Ensure repository has Actions enabled in Settings → Actions → General"
echo "2. Check that the workflow has proper permissions (contents: write)"
echo "3. Verify the tag is pushed to the remote repository"
echo "4. If using a fork, ensure Actions are enabled and you have write access"
echo "5. Consider using a Personal Access Token with 'repo' scope" 