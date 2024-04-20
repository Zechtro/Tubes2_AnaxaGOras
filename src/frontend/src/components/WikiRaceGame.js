import React, { useState } from 'react';
import Algorithm from './Algorithm';


const WikiRaceGame = () => {
  const [startPage, setStartPage] = useState('');
  const [targetPage, setTargetPage] = useState('');
  const [algorithm, setAlgorithm] = useState('');

  const handleAlgorithmChange = (event) => {
    setAlgorithm(event.target.value);
  };

  const handleStartGame = async () => {
    try {
      const result = await startWikiRace(startPage, targetPage, algorithm);
      console.log('Game Result:', result);      
    } catch (error) {
      console.error('Error starting game:', error);    
    }
  };

  return (
    <div className=''>
      <h2 className="text-3xl font-bold mb-6">Wiki Race Game</h2>

      <div className="flex flex-row mb-4">
        <div className="w-1/2 mr-2">
          <input
            type="text"
            placeholder="Start Page"
            value={startPage}
            onChange={(e) => setStartPage(e.target.value)}
            className="w-full px-4 py-3 bg-gray-700 border border-gray-600 rounded-lg focus:outline-none focus:ring focus:ring-blue-400 text-white"
          />
        </div>
        <div className="w-1/2 ml-2">
          <input
            type="text"
            placeholder="Target Page"
            value={targetPage}
            onChange={(e) => setTargetPage(e.target.value)}
            className="w-full px-4 py-3 bg-gray-700 border border-gray-600 rounded-lg focus:outline-none focus:ring focus:ring-blue-400 text-white"
          />
        </div>
      </div>

      <div className="flex justify-center mb-6 ml-5">
        <Algorithm
          label="BFS"
          value="bfs"
          checked={algorithm === 'bfs'}
          onChange={handleAlgorithmChange}
        />
        <Algorithm
          label="IDS"
          value="ids"
          checked={algorithm === 'ids'}
          onChange={handleAlgorithmChange}
        />
      </div>

      <button
        onClick={handleStartGame}
        className="w-full bg-blue-500 hover:bg-blue-600 text-white font-bold py-3 rounded-lg focus:outline-none focus:ring focus:ring-blue-400"
      >
        Start Game
      </button>
    </div>
  );
};

export default WikiRaceGame;
