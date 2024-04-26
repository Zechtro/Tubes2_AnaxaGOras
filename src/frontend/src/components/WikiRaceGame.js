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
        setResultTime(data.time)
      } else{
        setIsError(true)
        setErrorMsg(data.Error_Message)
      }
    } catch (error) {
      setIsError(true)
      setErrorMsg("Error Fetching Data")
      console.error('Error creating todo:', error);
    }finally {
      setIsLoading(false); // Set isLoading to false after fetch completes
    }
  };

  const graphz = {
    nodes: [
        {id: "1", label: "Node 1", title: "node 1 tooltip text", shape:"star",size: 15},
        {id: "2", label: "Timeline of the evolutionary history of life", title: "node 2 tooltip text",
        shape: "star", size:50,color: {
          border: "#222222",
          background: "#FFFFFF"
      },},
        {id: "3", label: "Node 3", title: "node 3 tooltip text",
        shape: "diamond"},
        {id: "4", label: "Node 4", title: "node 4 tooltip text",
        shape: "star"},
        {id: "5", label: "Node 5     ", title: "node 5 tooltip text"},
        {id: "6", label: "Snow", title: "node 6 tooltip text", shape: "circle"},
        {id: "7", label: "Node 7", title: "node 7 tooltip text"},
        {id: "8", label: "Node 8", title: "node 8 tooltip text"},
        {id: "9", label: "Node 9", title: "node 9 tooltip text",shape: "star",
        size: 15,
        color: {
            border: "#824D74",
            background: "#824D74"
        },
        font: {
            color: "#824D74",
            size: 15
        }}
    ],
    edges: [
        {from: "1", to: "1", smooth: {type: "curvedCW"}, arrows: {from: {enabled: true, type: "circle"}, to: {enabled: true, type: "circle"}}},
        {from: "1", to: "7", arrows: {from: {enabled: true, type: "vee"}, to: {enabled: true, type: "vee"}}},
        {from: "1", to: "3", arrows: {to: {enabled: true, type: "curve"}}},
        {from: "6", to: "5", color: {highlight: "#fff", opacity: 0.2}},
        {from: "6", to: "2"},
        {from: "7", to: "2"},
        {from: "6", to: "7"},
        {from: "6", to: "8"},
        {from: "7", to: "8"},
        {from: "8", to: "2"},
        {from: "3", to: "7"},
      ]
  }

  var options = {
    // groups:{
    //   group1 : {size:400}
    // },
    physics: {
        enabled: false
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
        {isResult && !isLoading && !isError && <p>Depth           : {resultDepth}</p>}
        {isResult && !isLoading && !isError && <p>Execution Time  : {resultTime} ms</p>}
      </div>
    </div>
  );
};

export default WikiRaceGame;
