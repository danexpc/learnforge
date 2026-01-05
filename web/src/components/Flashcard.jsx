import { useState } from 'preact/hooks';

export default function Flashcard({ flashcard, index }) {
  const [isFlipped, setIsFlipped] = useState(false);

  return (
    <div
      className="relative w-full h-48 cursor-pointer"
      style={{ perspective: '1000px' }}
      onClick={() => setIsFlipped(!isFlipped)}
    >
      <div
        className="relative w-full h-full transition-transform duration-500"
        style={{
          transformStyle: 'preserve-3d',
          transform: isFlipped ? 'rotateY(180deg)' : 'rotateY(0deg)',
        }}
      >
        {/* Front */}
        <div
          className="absolute w-full h-full bg-gradient-to-br from-blue-500 to-purple-600 rounded-xl shadow-lg p-6 flex items-center justify-center"
          style={{ backfaceVisibility: 'hidden', WebkitBackfaceVisibility: 'hidden' }}
        >
          <div className="text-center">
            <p className="text-white text-lg font-semibold mb-2">Question {index + 1}</p>
            <p className="text-white text-xl">{flashcard.q}</p>
            <p className="text-white/80 text-sm mt-4">Click to reveal answer</p>
          </div>
        </div>

        {/* Back */}
        <div
          className="absolute w-full h-full bg-gradient-to-br from-green-500 to-teal-600 rounded-xl shadow-lg p-6 flex items-center justify-center"
          style={{
            backfaceVisibility: 'hidden',
            WebkitBackfaceVisibility: 'hidden',
            transform: 'rotateY(180deg)',
          }}
        >
          <div className="text-center">
            <p className="text-white text-lg font-semibold mb-2">Answer</p>
            <p className="text-white text-xl">{flashcard.a}</p>
            <p className="text-white/80 text-sm mt-4">Click to flip back</p>
          </div>
        </div>
      </div>
    </div>
  );
}

