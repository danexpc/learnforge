import { useState } from 'preact/hooks';

export default function Meme({ memeUrl, topic, question, onRegenerate, regenerating }) {
  const [imageError, setImageError] = useState(false);

  return (
    <div className="bg-white rounded-lg shadow-md p-3 border border-gray-200">
      <div className="flex items-center justify-between mb-2">
        <h3 className="text-xs font-semibold text-gray-600">Meme</h3>
        <button
          onClick={onRegenerate}
          disabled={regenerating}
          className="w-6 h-6 flex items-center justify-center rounded bg-gray-100 hover:bg-gray-200 text-gray-700 hover:text-gray-900 disabled:opacity-50 disabled:cursor-not-allowed transition-colors text-xs"
          title="Regenerate meme"
        >
          {regenerating ? (
            <div className="w-3 h-3 border-2 border-gray-400 border-t-transparent rounded-full animate-spin"></div>
          ) : (
            '🔄'
          )}
        </button>
      </div>
      
      {regenerating && !memeUrl ? (
        <div className="flex items-center justify-center h-24 bg-gray-50 rounded">
          <div className="text-center">
            <div className="inline-block w-4 h-4 border-2 border-blue-600 border-t-transparent rounded-full animate-spin mb-1"></div>
            <p className="text-xs text-gray-500">Generating...</p>
          </div>
        </div>
      ) : memeUrl && !imageError ? (
        <div className="flex justify-center">
          <img
            key={memeUrl} // Force re-render on URL change to prevent duplicates
            src={memeUrl}
            alt={`Meme about ${topic}`}
            className="w-full h-auto rounded-lg shadow-sm cursor-pointer hover:shadow-md transition-shadow"
            onError={() => setImageError(true)}
            onClick={() => window.open(memeUrl, '_blank')}
            title="Click to view full size"
          />
        </div>
      ) : imageError ? (
        <div className="text-center py-3 text-xs text-gray-500">
          Failed to load
        </div>
      ) : (
        <div className="flex items-center justify-center h-24 bg-gray-50 rounded">
          <p className="text-xs text-gray-400">No meme</p>
        </div>
      )}
    </div>
  );
}

