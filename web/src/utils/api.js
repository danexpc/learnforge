const API_BASE = '';

export async function processText(data) {
  const response = await fetch(`${API_BASE}/v1/process`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(data),
  });

  if (!response.ok) {
    const error = await response.json();
    throw new Error(error.error?.message || 'Failed to process text');
  }

  return response.json();
}

export async function getResult(id) {
  const response = await fetch(`${API_BASE}/v1/process/${id}`);

  if (!response.ok) {
    const error = await response.json();
    throw new Error(error.error?.message || 'Failed to get result');
  }

  return response.json();
}

export async function regenerateMeme(topic, question) {
  const response = await fetch(`${API_BASE}/v1/process`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({
      text: `${topic}. ${question || 'Educational content'}`,
      mode: 'flashcards',
      topic: topic,
      generate_meme: true,
    }),
  });

  if (!response.ok) {
    const error = await response.json();
    throw new Error(error.error?.message || 'Failed to regenerate meme');
  }

  const data = await response.json();
  return data.meme_url;
}

