<template>
  <div class="card">
    <a :href="props.url" target="_blank" rel="noopener" class="card-link">
      <div class="card-header">
        <img :src="props.author.avatarURL" alt="Author Avatar" class="avatar" />
        <span class="author-name">{{ props.author.username }}</span>
        <span class="time">{{ new Date(props.time).toLocaleString() }}</span>
      </div>
      <div class="card-message">
        {{ props.message }}
      </div>
    </a>
  </div>
</template>

<script setup lang="ts">
const props = defineProps<{
  author: {
    username: string;
    avatarURL: string;
  };
  message: string;
  time: string;
  url: string;
}>();
</script>

<style scoped>
.card {
  background: white;
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  overflow: hidden;
  transition: box-shadow 0.3s ease;
}

.card:hover {
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.15);
}

.card-link {
  display: block;
  text-decoration: none;
  color: inherit;
}

.card-header {
  display: flex;
  align-items: center;
  padding: 16px;
  border-bottom: 1px solid #e5e7eb;
  background-color: #f9fafb;
  gap: 12px;
}

.avatar {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  object-fit: cover;
  border: 2px solid #e5e7eb;
}

.author-name {
  font-weight: 600;
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', 'Roboto', 'Helvetica Neue', Arial, sans-serif;
  font-size: 0.95rem;
  color: #1f2937;
  flex: 1;
  letter-spacing: -0.01em;
}

.time {
  font-size: 0.85rem;
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', 'Roboto', 'Helvetica Neue', Arial, sans-serif;
  font-weight: 500;
  color: #6b7280;
  white-space: nowrap;
  letter-spacing: -0.005em;
}

.card-message {
  padding: 16px;
  color: #ffffff;
  line-height: 1.5;
  font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
  background-color: #000000;
  white-space: pre-wrap;
  word-break: break-word;
  border-radius: 4px;
  margin: 8px;
}
</style>