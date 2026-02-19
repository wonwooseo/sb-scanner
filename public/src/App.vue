<template>
  <div class="app-container">
    <header class="app-header">
      <h1>
        <a href="https://github.com/wonwooseo/sb-scanner" target="_blank" rel="noopener noreferrer">SB Scanner</a>
      </h1>
      <p>커밋 메시지를 예쁘게 씁시다.</p>
    </header>
    <div class="cards-container">
      <Card
        v-for="commit in commits"
        :key="commit.sha"
        :author="{
          username: commit.author.username,
          avatarURL: commit.author.avatar_url,
        }"
        :message="commit.message"
        :time="commit.time"
        :sentiment="{
          score: commit.sentiment.score,
          model: commit.sentiment.model,
        }"
        :url="commit.url"
      />
      <div v-if="loading" class="loading-indicator">
        <div class="spinner"></div>
        <span>Loading more commits...</span>
      </div>
      <div v-if="!hasMore && commits.length > 0" class="end-message">
        No more commits to load
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import Card from '@/components/Card.vue';

import { onMounted, onUnmounted, ref } from 'vue';
import { api } from '@/api/api';
import type { Commit } from '@/api/api';

const commits = ref<Commit[]>([]);
const loading = ref(false);
const hasMore = ref(true);
const currentBookmark = ref<string | null>(null);

const getCommits = async () => {
  try {
    loading.value = true;
    const data = await api.getCommits(null, 10);
    commits.value = data.commits;
    currentBookmark.value = data.bookmark || null;
    hasMore.value = !!data.bookmark;
  } catch (error) {
    console.error('Error fetching commits:', error);
  } finally {
    loading.value = false;
  }
};

const loadMoreCommits = async () => {
  if (loading.value || !hasMore.value) return;
  
  try {
    loading.value = true;
    const data = await api.getCommits(currentBookmark.value, 10);
    commits.value.push(...data.commits);
    currentBookmark.value = data.bookmark || null;
    hasMore.value = !!data.bookmark;
  } catch (error) {
    console.error('Error loading more commits:', error);
  } finally {
    loading.value = false;
  }
};

const handleScroll = () => {
  const { scrollTop, scrollHeight, clientHeight } = document.documentElement;
  // Load more when user is 100px from bottom
  if (scrollTop + clientHeight >= scrollHeight - 100) {
    loadMoreCommits();
  }
};

onMounted(() => {
  getCommits();
  window.addEventListener('scroll', handleScroll);
});

onUnmounted(() => {
  window.removeEventListener('scroll', handleScroll);
});
</script>

<style>
body {
    font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', 'Roboto', 'Helvetica Neue', Arial, sans-serif;
}

.app-container {
  min-height: 100vh;
  background-color: #ffffff;
  padding: 20px;
  display: flex;
  flex-direction: column;
  align-items: center;
}

.app-header {
  text-align: center;
  margin-bottom: 32px;
}

.app-header h1 a {
  font-size: 3rem;
  font-weight: 700;
  color: #1f2937;
  margin: 0 0 8px 0;
  text-decoration: none;
}

.app-header h1 a:hover {
  color: #3b82f6;
  text-decoration: underline;
}

.app-header p {
  font-size: 1.2rem;
  color: #6b7280;
  margin: 0;
  font-weight: 400;
}

.cards-container {
  max-width: 600px;
  width: 100%;
  display: flex;
  flex-direction: column;
  gap: 24px;
  background-color: #f1f5f9;
  padding: 24px;
  border-radius: 12px;
  box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1), 0 2px 4px -1px rgba(0, 0, 0, 0.06);
  border: 1px solid #e2e8f0;
}

.loading-indicator {
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 24px;
  gap: 12px;
  color: #6b7280;
  font-size: 0.9rem;
}

.spinner {
  width: 20px;
  height: 20px;
  border: 2px solid #e5e7eb;
  border-top: 2px solid #3b82f6;
  border-radius: 50%;
  animation: spin 1s linear infinite;
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}

.end-message {
  text-align: center;
  padding: 24px;
  color: #6b7280;
  font-size: 0.9rem;
  font-style: italic;
}

@media (max-width: 640px) {
  .app-container {
    padding: 12px;
  }

  .app-header h1 {
    font-size: 1.75rem;
  }

  .app-header p {
    font-size: 0.95rem;
  }

  .cards-container {
    padding: 16px;
    gap: 16px;
  }
}
</style>
