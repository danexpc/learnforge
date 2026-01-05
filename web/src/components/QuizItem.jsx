import { useState } from 'preact/hooks';

export default function QuizItem({ quiz, index }) {
  const [selectedAnswer, setSelectedAnswer] = useState(null);
  const [showResult, setShowResult] = useState(false);

  const handleSelect = (choice) => {
    if (showResult) return;
    // Store the original choice for comparison
    setSelectedAnswer(choice);
    setShowResult(true);
  };

  // Helper to clean answer strings (remove letter prefixes)
  const cleanAnswer = (str) => str ? str.replace(/^[A-Z]\.\s*/, '').trim() : '';
  
  // Compare answers, handling both with and without letter prefixes
  const cleanQuizAnswer = cleanAnswer(quiz.answer);
  const cleanSelected = cleanAnswer(selectedAnswer);
  const isCorrect = selectedAnswer === quiz.answer || 
                    cleanSelected === cleanQuizAnswer ||
                    cleanSelected === quiz.answer ||
                    selectedAnswer === cleanQuizAnswer;

  return (
    <div className="bg-white rounded-xl shadow-lg p-6 border border-gray-100 mb-4">
      <div className="flex items-start mb-4">
        <span className="flex-shrink-0 w-8 h-8 bg-blue-600 text-white rounded-full flex items-center justify-center font-bold mr-3">
          {index + 1}
        </span>
        <h4 className="text-lg font-semibold text-gray-800 flex-1">{quiz.q}</h4>
      </div>

      <div className="space-y-2 ml-11">
        {quiz.choices.map((choice, idx) => {
          // Strip any existing letter prefix (A., B., etc.) from the choice
          const cleanChoice = choice.replace(/^[A-Z]\.\s*/, '');
          const letter = String.fromCharCode(65 + idx);
          
          let buttonClass = 'w-full text-left px-4 py-3 rounded-lg border-2 transition-all duration-200 ';
          
          if (showResult) {
            // Compare using both original and cleaned versions
            const cleanQuizAnswer = cleanAnswer(quiz.answer);
            const cleanSelectedAnswer = cleanAnswer(selectedAnswer);
            const isCorrectAnswer = choice === quiz.answer || cleanChoice === quiz.answer || cleanChoice === cleanQuizAnswer;
            const isSelected = choice === selectedAnswer || cleanChoice === cleanSelectedAnswer;
            
            if (isCorrectAnswer) {
              buttonClass += 'bg-green-100 border-green-500 text-green-800 font-semibold';
            } else if (isSelected && !isCorrectAnswer) {
              buttonClass += 'bg-red-100 border-red-500 text-red-800';
            } else {
              buttonClass += 'bg-gray-50 border-gray-200 text-gray-600';
            }
          } else {
            buttonClass += 'bg-white border-gray-300 hover:border-blue-500 hover:bg-blue-50 text-gray-700 cursor-pointer';
          }

          return (
            <button
              key={idx}
              onClick={() => handleSelect(choice)}
              disabled={showResult}
              className={buttonClass}
            >
              <span className="font-medium">{letter}.</span> {cleanChoice}
            </button>
          );
        })}
      </div>

      {showResult && (
        <div className={`mt-4 ml-11 p-3 rounded-lg ${isCorrect ? 'bg-green-50 text-green-800' : 'bg-red-50 text-red-800'}`}>
          {isCorrect ? (
            <p className="font-semibold">✓ Correct! Well done.</p>
          ) : (
            <p className="font-semibold">✗ Incorrect. The correct answer is: <strong>{cleanAnswer(quiz.answer) || quiz.answer}</strong></p>
          )}
        </div>
      )}
    </div>
  );
}

