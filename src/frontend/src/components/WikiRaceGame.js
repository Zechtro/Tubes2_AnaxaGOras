import React, { useState } from 'react';
import Algorithm from './Algorithm';
import Graph from 'react-vis-network-graph'


const WikiRaceGame = () => {
  const [StartPage, setStartPage] = useState('');
  const [TargetPage, setTargetPage] = useState('');
  const [AlgorithmUsed, setAlgorithm] = useState('');

  const [isLoading, setIsLoading] = useState(false)

  const [isResult, setIsResult] = useState(false)
  const [resultGraph, setResultGraph] = useState({nodes:[],edges:[]})
  const [resultDepth, setResultDepth] = useState(0)
  const [articleChecked, setArticleChecked] = useState(0)
  const [resultTime, setResultTime] = useState(0)
  
  const [isError, setIsError] = useState(false)
  const [errorMsg, setErrorMsg] = useState('')

  const handleAlgorithmChange = (event) => {
    setAlgorithm(event.target.value);
  };

  const handleStartGame = async () => {
    // try {
    //   const result = await startWikiRace(startPage, targetPage, algorithm);
    //   console.log('Game Result:', result);      
    // } catch (error) {
    //   console.error('Error starting game:', error);    
    // }
    if (StartPage === "" || TargetPage === ""){
      setIsError(true)
      setResultGraph({nodes:[],edges:[]})
      setErrorMsg("Start Page and Target Page must not empty")
    }else if (StartPage === TargetPage){
      setIsError(true)
      setResultGraph({nodes:[],edges:[]})
      setErrorMsg("Start Page must not be the same with Target Page")
    }else{
      setIsLoading(true)
      try {
        const response = await fetch('http://localhost:8000/api/process', {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({algorithm: AlgorithmUsed, startPage: StartPage, targetPage : TargetPage})
        });
    
        if (!response.ok) {
          throw new Error('Network response was not ok');
        }
    
        const data = await response.json(); 
        if (data.status  === "OK"){
          setIsError(false)
          setIsResult(true)
          setResultGraph(data.graph)
          setResultDepth(data.depth)
          setArticleChecked(data.checked)
          setResultTime(data.time)
        } else{
          setResultGraph({nodes:[],edges:[]})
          setIsError(true)
          setErrorMsg(data.Error_Message)
        }
      } catch (error) {
        setResultGraph({nodes:[],edges:[]})
        setIsError(true)
        setErrorMsg("Error Fetching Data")
        console.error('Error creating todo:', error);
      }finally {
        setIsLoading(false); // Set isLoading to false after fetch completes
      }
    }
  };

  var options = {
    physics: {
        enabled: true
    },
    interaction: {
        navigationButtons: true
    },
    edges: {
        color: "white"
    },
    shadow: true,
    smooth: true,
    height: "500px",
    innerWidth:"400px"
  }

  return (
    <div className=''>
      <div className=''>
        <h2 className="text-3xl font-bold mb-6">Wiki Race Game</h2>

        <div className="flex flex-row mb-4">
          <div className="w-1/2 mr-2">
            <input
              type="text"
              placeholder="Start Page"
              value={StartPage}
              onChange={(e) => setStartPage(e.target.value)}
              className="w-full px-4 py-3 bg-gray-700 border border-gray-600 rounded-lg focus:outline-none focus:ring focus:ring-blue-400 text-white"
            />
          </div>
          <div className="w-1/2 ml-2">
            <input
              type="text"
              placeholder="Target Page"
              value={TargetPage}
              onChange={(e) => setTargetPage(e.target.value)}
              className="w-full px-4 py-3 bg-gray-700 border border-gray-600 rounded-lg focus:outline-none focus:ring focus:ring-blue-400 text-white"
            />
          </div>
        </div>

        <div className="flex justify-around mb-6">
          <Algorithm
            label="BFS"
            value="bfs"
            checked={AlgorithmUsed === 'bfs'}
            onChange={handleAlgorithmChange}
          />
          <Algorithm
            label="IDS"
            value="ids"
            checked={AlgorithmUsed === 'ids'}
            onChange={handleAlgorithmChange}
          />
        </div>

        {!isLoading && <button
          onClick={handleStartGame}
          className="w-full bg-blue-500 hover:bg-blue-600 text-white font-bold py-3 rounded-lg focus:outline-none focus:ring focus:ring-blue-400"
        >
          Start Game
        </button>}
      </div>
      <div className="w-400 flex items-center justify-center">
        {isLoading && <p>Processing...</p>}
        {isError && !isLoading && <p>{errorMsg}</p>}
        {isResult && !isLoading && !isError && <Graph graph={resultGraph} options={options}/>}
      </div>
      <div>
        {isResult && !isLoading && !isError && <p>Depth                 : {resultDepth}</p>}
        {isResult && !isLoading && !isError && <p>Total Checked Article : {articleChecked}</p>}
        {isResult && !isLoading && !isError && <p>Execution Time        : {resultTime} ms</p>}
      </div>
    </div>
  );
};

export default WikiRaceGame;
