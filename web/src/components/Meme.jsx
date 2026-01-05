import { useState, useEffect } from 'preact/hooks';

export default function Meme({ memeUrl, topic, question, onRegenerate, regenerating, expectingMeme = false }) {
  const [imageError, setImageError] = useState(false);
  const [imageLoading, setImageLoading] = useState(!!memeUrl);

  useEffect(() => {
    if (memeUrl) {
      setImageLoading(true);
      setImageError(false);
    }
  }, [memeUrl]);

  const contentHeight = "min-h-[200px]";

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
      
      <div className={`${contentHeight} relative bg-gray-50 rounded-lg overflow-hidden flex items-center justify-center`}>
        {regenerating && !memeUrl ? (
          <div className="text-center">
            <div className="inline-block w-6 h-6 border-2 border-blue-600 border-t-transparent rounded-full animate-spin mb-2"></div>
            <p className="text-xs text-gray-500">Generating...</p>
          </div>
        ) : expectingMeme && !memeUrl ? (
          <div className="text-center">
            <div className="inline-block w-6 h-6 border-2 border-blue-600 border-t-transparent rounded-full animate-spin mb-2"></div>
            <p className="text-xs text-gray-500">Generating...</p>
          </div>
        ) : memeUrl ? (
          <>
            <div className={`absolute inset-0 flex items-center justify-center z-10 transition-opacity duration-200 ${
              imageLoading ? 'opacity-100' : 'opacity-0 pointer-events-none'
            }`}>
              <div className="text-center">
                <div className="inline-block w-6 h-6 border-2 border-blue-600 border-t-transparent rounded-full animate-spin mb-2"></div>
                <p className="text-xs text-gray-500">Loading...</p>
              </div>
            </div>
            <img
              key={memeUrl}
              src={memeUrl}
              alt={`Meme about ${topic}`}
              className={`w-full h-auto rounded-lg shadow-sm cursor-pointer hover:shadow-md transition-opacity duration-300 ${
                imageLoading ? 'opacity-0' : 'opacity-100'
              }`}
              onLoad={() => setImageLoading(false)}
              onError={() => {
                setImageError(true);
                setImageLoading(false);
              }}
              onClick={() => window.open(memeUrl, '_blank')}
              title="Click to view full size"
            />
          </>
        ) : imageError ? (
          <div className="text-center text-xs text-gray-500">
            Failed to load
          </div>
        ) : (
          <p className="text-xs text-gray-400">No meme</p>
        )}
      </div>
    </div>
  );
}

