import { useState, useEffect } from 'preact/hooks';
import { processText, regenerateMeme, getResult } from './utils/api';
import { saveLesson } from './utils/storage';
import Button from './components/Button';
import Card from './components/Card';
import Loading from './components/Loading';
import Flashcard from './components/Flashcard';
import QuizItem from './components/QuizItem';
import Meme from './components/Meme';
import SavedLessons from './components/SavedLessons';

export default function App() {
  const [text, setText] = useState('');
  const [mode, setMode] = useState('lesson');
  const [topic, setTopic] = useState('');
  const [level, setLevel] = useState('beginner');
  const [generateMeme, setGenerateMeme] = useState(false);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);
  const [result, setResult] = useState(null);
  const [regeneratingMeme, setRegeneratingMeme] = useState({});
  const [memeUrls, setMemeUrls] = useState([]);

  const generateAdditionalMemes = async (topic, question) => {
    for (let i = 0; i < 2; i++) {
      try {
        const newMemeUrl = await regenerateMeme(topic, question);
        if (newMemeUrl) {
          setMemeUrls(prev => [...prev, newMemeUrl]);
        }
      } catch (err) {
        console.error('Failed to generate additional meme:', err);
      }
    }
  };

  const loadLessonById = async (id) => {
    setLoading(true);
    setError(null);
    try {
      const response = await getResult(id);
      setResult(response);
      
      if (response.topic) {
        setTopic(response.topic);
      }
      
      if (response.meme_url) {
        setMemeUrls([response.meme_url]);
        setGenerateMeme(true);
        generateAdditionalMemes(response.topic, response.quiz?.[0]?.q || response.flashcards?.[0]?.q);
      } else {
        setMemeUrls([]);
      }

      const newUrl = new URL(window.location);
      newUrl.searchParams.set('id', id);
      window.history.pushState({}, '', newUrl);
    } catch (err) {
      if (err.message.includes('not found') || err.message.includes('404')) {
        setError('This lesson no longer exists. It may have been deleted or expired.');
      } else {
        setError(err.message || 'Failed to load content');
      }
      const newUrl = new URL(window.location);
      newUrl.searchParams.delete('id');
      window.history.replaceState({}, '', newUrl);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    const params = new URLSearchParams(window.location.search);
    const id = params.get('id');
    
    if (id) {
      loadLessonById(id);
    }
  }, []);

  const handleSubmit = async (e) => {
    e.preventDefault();
    setLoading(true);
    setError(null);
    setResult(null);
    setMemeUrls([]);

    try {
      const data = {
        text,
        mode,
        ...(topic && { topic }),
        level,
        generate_meme: generateMeme,
      };

      const response = await processText(data);
      setResult(response);
      
      if (response.id) {
        const newUrl = new URL(window.location);
        newUrl.searchParams.set('id', response.id);
        window.history.pushState({}, '', newUrl);

        saveLesson({
          id: response.id,
          topic: response.topic,
          mode: mode,
          level: level,
          createdAt: response.created_at || new Date().toISOString(),
        });
      }
      
      if (response.meme_url && generateMeme) {
        setMemeUrls([response.meme_url]);
        generateAdditionalMemes(response.topic, response.quiz?.[0]?.q || response.flashcards?.[0]?.q);
      } else {
        setMemeUrls([]);
      }
    } catch (err) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  };

  const handleReset = () => {
    setText('');
    setMode('lesson');
    setTopic('');
    setLevel('beginner');
    setGenerateMeme(false);
    setResult(null);
    setError(null);
      setRegeneratingMeme({});
      setMemeUrls([]);
      
      const newUrl = new URL(window.location);
    newUrl.searchParams.delete('id');
    window.history.replaceState({}, '', newUrl);
  };

  const handleRegenerateMeme = async (index) => {
    if (!result) return;
    
    setRegeneratingMeme(prev => ({ ...prev, [index]: true }));
    try {
      // Get question for meme
      const question = result.quiz?.[0]?.q || result.flashcards?.[0]?.q || '';
      
      const newMemeUrl = await regenerateMeme(result.topic, question);
      if (newMemeUrl) {
        setMemeUrls(prev => {
          const updated = [...prev];
          updated[index] = newMemeUrl;
          return updated;
        });
      }
    } catch (err) {
      console.error('Failed to regenerate meme:', err);
      // Silently fail - meme generation is optional
    } finally {
      setRegeneratingMeme(prev => ({ ...prev, [index]: false }));
    }
  };

  return (
    <div className="min-h-screen bg-gradient-to-br from-blue-50 via-white to-purple-50 flex flex-col">
      {/* Header */}
      <header className="bg-white shadow-sm border-b border-gray-200">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-6">
          <div className="flex items-center justify-between">
            <div>
              <h1 className="text-3xl font-bold bg-gradient-to-r from-blue-600 to-purple-600 bg-clip-text text-transparent">
                LearnForge
              </h1>
              <p className="text-gray-600 mt-1">AI-Powered Learning Content Generator</p>
            </div>
            <a
              href="/docs"
              target="_blank"
              className="text-blue-600 hover:text-blue-700 font-medium"
            >
              API Docs →
            </a>
          </div>
        </div>
      </header>

      {/* Main Content */}
      <main className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8 flex-1">
        {!result ? (
          <div className="max-w-3xl mx-auto space-y-6">
            <Card title="Transform Text into Learning Content">
              <div className="relative">
                <form onSubmit={handleSubmit} className={`space-y-6 ${loading ? 'opacity-50 pointer-events-none' : ''}`}>
                {/* Text Input */}
                <div>
                  <label htmlFor="text" className="block text-sm font-medium text-gray-700 mb-2">
                    Enter your text
                  </label>
                  <textarea
                    id="text"
                    value={text}
                    onChange={(e) => setText(e.target.value)}
                    required
                    rows="8"
                    className="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500 resize-none"
                    placeholder="Paste or type your text here. For example: 'Photosynthesis is the process by which plants convert light energy into chemical energy...'"
                  />
                </div>

                {/* Options */}
                <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
                  <div>
                    <label htmlFor="mode" className="block text-sm font-medium text-gray-700 mb-2">
                      Mode
                    </label>
                    <select
                      id="mode"
                      value={mode}
                      onChange={(e) => setMode(e.target.value)}
                      className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                    >
                      <option value="lesson">Lesson</option>
                      <option value="flashcards">Flashcards</option>
                      <option value="quiz">Quiz</option>
                    </select>
                  </div>

                  <div>
                    <label htmlFor="level" className="block text-sm font-medium text-gray-700 mb-2">
                      Level
                    </label>
                    <select
                      id="level"
                      value={level}
                      onChange={(e) => setLevel(e.target.value)}
                      className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                    >
                      <option value="beginner">Beginner</option>
                      <option value="intermediate">Intermediate</option>
                      <option value="advanced">Advanced</option>
                    </select>
                  </div>

                  <div>
                    <label htmlFor="topic" className="block text-sm font-medium text-gray-700 mb-2">
                      Topic (optional)
                    </label>
                    <input
                      id="topic"
                      type="text"
                      value={topic}
                      onChange={(e) => setTopic(e.target.value)}
                      placeholder="e.g., Science, History"
                      className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                    />
                  </div>
                </div>

                {/* Generate Meme Option */}
                <div className="flex items-center">
                  <input
                    id="generateMeme"
                    type="checkbox"
                    checked={generateMeme}
                    onChange={(e) => setGenerateMeme(e.target.checked)}
                    className="w-4 h-4 text-blue-600 border-gray-300 rounded focus:ring-blue-500"
                  />
                  <label htmlFor="generateMeme" className="ml-2 text-sm text-gray-700">
                    Generate memes related to the topic <span className="text-xs bg-yellow-100 text-yellow-800 px-2 py-0.5 rounded ml-1">BETA</span>
                  </label>
                </div>

                {/* Error Message */}
                {error && (
                  <div className="bg-red-50 border border-red-200 rounded-lg p-4">
                    <p className="text-red-800 font-medium">Error: {error}</p>
                  </div>
                )}

                {/* Submit Button */}
                <Button
                  type="submit"
                  disabled={loading || !text.trim()}
                  variant="primary"
                  className="w-full"
                >
                  {loading ? 'Processing...' : 'Generate Learning Content'}
                </Button>
              </form>
              
              {loading && (
                <div className="absolute inset-0 bg-white bg-opacity-95 flex items-center justify-center rounded-lg z-10">
                  <Loading message="AI is processing your text and generating learning content..." />
                </div>
              )}
              </div>
            </Card>

            {/* Saved Lessons - Below form */}
            <SavedLessons onLoadLesson={loadLessonById} />
          </div>
        ) : (
          <div className="space-y-6">
            {/* Result Header */}
            <div className="flex items-start justify-between gap-4">
              <div className="flex-1">
                <h2 className="text-2xl font-bold text-gray-800">Generated Content</h2>
                <p className="text-gray-600 mt-1">
                  Topic: <span className="font-semibold">{result.topic}</span> • 
                  Level: <span className="font-semibold capitalize">{level}</span>
                </p>
              </div>
              <div className="flex flex-col gap-2 items-end">
                <Button onClick={handleReset} variant="secondary">
                  Create New
                </Button>
              </div>
            </div>

            {/* Memes - 3 in a row */}
            {generateMeme && result && (
              <Card title="🎨 Memes (Beta)">
                <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
                  {[0, 1, 2].map((index) => (
                    <Meme
                      key={index}
                      memeUrl={memeUrls[index]}
                      topic={result.topic}
                      question={result.quiz?.[0]?.q || result.flashcards?.[0]?.q}
                      onRegenerate={() => handleRegenerateMeme(index)}
                      regenerating={regeneratingMeme[index] || false}
                      expectingMeme={true}
                    />
                  ))}
                </div>
              </Card>
            )}

            {/* Summary */}
            {result.summary && (
              <Card title="Summary">
                <p className="text-gray-700 leading-relaxed whitespace-pre-wrap">{result.summary}</p>
              </Card>
            )}

            {/* Key Points */}
            {result.key_points && result.key_points.length > 0 && (
              <Card title="Key Points">
                <ul className="space-y-2">
                  {result.key_points.map((point, idx) => (
                    <li key={idx} className="flex items-start">
                      <span className="flex-shrink-0 w-6 h-6 bg-blue-100 text-blue-600 rounded-full flex items-center justify-center font-semibold mr-3 mt-0.5">
                        {idx + 1}
                      </span>
                      <span className="text-gray-700 flex-1">{point}</span>
                    </li>
                  ))}
                </ul>
              </Card>
            )}

            {/* Flashcards */}
            {result.flashcards && result.flashcards.length > 0 && (
              <Card title="Flashcards">
                <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                  {result.flashcards.map((flashcard, idx) => (
                    <Flashcard key={idx} flashcard={flashcard} index={idx} />
                  ))}
                </div>
              </Card>
            )}

            {/* Quiz */}
            {result.quiz && result.quiz.length > 0 && (
              <Card title="Quiz Questions">
                {result.quiz.map((quiz, idx) => (
                  <QuizItem key={idx} quiz={quiz} index={idx} />
                ))}
              </Card>
            )}

            {/* Metadata */}
            <Card className="bg-gray-50">
              <div className="text-sm text-gray-600 space-y-1">
                <p>Generated in {result.meta?.processing_ms}ms</p>
                <p>Model: {result.meta?.model} ({result.meta?.provider})</p>
                <p>ID: <code className="bg-gray-200 px-2 py-1 rounded">{result.id}</code></p>
              </div>
            </Card>
          </div>
        )}
      </main>

      {/* Footer */}
      <footer className="bg-white border-t border-gray-200 mt-auto">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-6">
          <p className="text-center text-gray-600 text-sm">
            Powered by AI • Built with Preact & Go
          </p>
        </div>
      </footer>
    </div>
  );
}

