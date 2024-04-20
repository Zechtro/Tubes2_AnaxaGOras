import { Routes, Route } from "react-router-dom";
import './App.css';
import WikiRaceGame from "./components/WikiRaceGame";

function App() {
  return (
    <div className="text-white mx-auto my-auto w-full p-2 lg:p-12 text-center">
      <WikiRaceGame />
    </div>
  );
}

export default App;
