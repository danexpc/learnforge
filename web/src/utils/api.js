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

