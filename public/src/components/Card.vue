<template>
  <div class="card">
    <a :href="props.url" target="_blank" rel="noopener" class="card-link">
      <div class="card-header">
        <img :src="props.author.avatarURL || 'https://github.githubassets.com/images/gravatars/gravatar-user-420.png?size=40'" alt="Author Avatar" class="avatar" />
        <span class="author-name">{{ props.author.username || 'user' }}</span>
        <span class="time">{{ new Date(props.time).toLocaleString() }}</span>
        <span class="sentiment" :class="`sentiment-${getSentimentLabel()}`" :data-tooltip="`Sentiment: ${props.sentiment.score.toFixed(2)}\nModel: ${props.sentiment.model}`">
          {{ getSentimentLabel().toUpperCase() }}
        </span>
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
  sentiment: {
    score: number;
    model: string;
  }
  url: string;
}>();

const getSentimentLabel = () => {
  const score = props.sentiment.score;
  if (score >= 0.3) return 'positive';
  if (score <= -0.3) return 'negative';
  return 'neutral';
};
</script>

<style scoped>
.card {
  background: white;
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  overflow: visible;
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
  font-size: 0.95rem;
  color: #1f2937;
  flex: 1;
  letter-spacing: -0.01em;
}

.time {
  font-size: 0.85rem;
  font-weight: 500;
  color: #6b7280;
  white-space: nowrap;
  letter-spacing: -0.005em;
}

.sentiment {
  font-size: 0.85rem;
  font-weight: 600;
  padding: 4px 8px;
  border-radius: 4px;
  white-space: nowrap;
  position: relative;
  cursor: pointer;
}

.sentiment::after {
  content: attr(data-tooltip);
  position: absolute;
  bottom: 125%;
  left: 50%;
  transform: translateX(-50%);
  background: #1f2937;
  color: white;
  padding: 8px;
  border-radius: 4px;
  font-size: 0.8rem;
  width: max-content;
  white-space: pre-line;
  opacity: 0;
  pointer-events: none;
  transition: opacity 0.2s;
  z-index: 10;
}

.sentiment:hover::after {
  opacity: 1;
}

@media (max-width: 640px) {
  .card-header {
    padding: 12px;
    gap: 8px;
  }

  .author-name {
    flex: 1 1 auto;
    min-width: 0;
    font-size: 0.9rem;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }

  .time {
    flex: 0 0 auto;
    font-size: 0.75rem;
    white-space: nowrap;
  }

  .sentiment {
    flex: 0 0 auto;
    font-size: 0.75rem;
    padding: 3px 6px;
  }

  .sentiment::after {
    width: 140px;
    font-size: 0.7rem;
    padding: 6px;
    bottom: auto;
    top: 125%;
  }

  .card-message {
    padding: 12px;
    margin: 4px;
  }
}

.sentiment-positive {
  background-color: #d1fae5;
  color: #065f46;
}

.sentiment-neutral {
  background-color: #f3f4f6;
  color: #374151;
}

.sentiment-negative {
  background-color: #fee2e2;
  color: #7f1d1d;
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