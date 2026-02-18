import axios from 'axios';

const apiClient = axios.create({
  baseURL: '/api/v1',
  headers: {
    'Content-Type': 'application/json',
  },
});

export const api = {
  async getCommits(bookmark: string | null, limit: number) : Promise<SearchCommitResponse> {
    const params: any = { limit };
    if (bookmark) {
      params.bookmark = bookmark;
    }
    const resp = await apiClient.get('/commit', { params });
    if (resp.status !== 200) {
        throw new Error(`HTTP error! status: ${resp.status}`);
    }
    return resp.data;
  },
};

export interface SearchCommitResponse {
  commits: Commit[];
  bookmark?: string;
}

export interface Commit {
  sha: string;
  url: string;
  message: string;
  author: Author;
  time: string;
  sentiment: Sentiment;
}

export interface Author {
  username: string;
  avatar_url: string;
}

export interface Sentiment {
  score: number;
  model: string;
}
