import React, { useEffect, useState } from 'react';
import './App.css';
import { Lume, LumeType } from './genproto/lume/v1/lume_pb.js';
import { createLume } from './clients/LumeClient';
import { LumeCard } from './components/LumeCard';

function App() {
  const [lumes, setLumes] = useState<Lume[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const createSampleLumes = async () => {
      try {
        setLoading(true);

        // Create sample lumes
        const sampleLumes = [
          {
            name: "San Francisco",
            description: "Golden Gate City",
            latitude: 37.7749,
            longitude: -122.4194,
            address: "San Francisco, CA, USA",
            type: LumeType.CITY
          },
          {
            name: "Golden Gate Bridge", 
            description: "Iconic suspension bridge",
            latitude: 37.8199,
            longitude: -122.4783,
            address: "Golden Gate Bridge, San Francisco, CA",
            type: LumeType.ATTRACTION
          },
          {
            name: "Fisherman's Wharf",
            description: "Popular tourist destination",
            latitude: 37.8080,
            longitude: -122.4177,
            address: "Fisherman's Wharf, San Francisco, CA",
            type: LumeType.ATTRACTION
          }
        ];

        const createdLumes: Lume[] = [];
        
        for (const lumeData of sampleLumes) {
          const response = await createLume(
            lumeData.name,
            lumeData.description,
            lumeData.latitude,
            lumeData.longitude,
            lumeData.address,
            lumeData.type
          );
          
          if (response.lume) {
            createdLumes.push(response.lume);
          }
        }

        setLumes(createdLumes);
        setError(null);
      } catch (err) {
        console.error("Error creating lumes:", err);
        setError("Failed to create lumes. Make sure the backend is running.");
      } finally {
        setLoading(false);
      }
    };

    createSampleLumes();
  }, []);

  if (loading) {
    return (
      <div className="App">
        <header className="App-header">
          <h1>Lumes</h1>
          <p>Simple Lume Management</p>
        </header>
        <main className="App-content">
          <p>Creating Lumes...</p>
        </main>
      </div>
    );
  }

  if (error) {
    return (
      <div className="App">
        <header className="App-header">
          <h1>Lumes</h1>
          <p>Simple Lume Management</p>
        </header>
        <main className="App-content">
          <div className="error-message">
            <p>{error}</p>
            <p>Make sure the Lume service is running on the backend.</p>
          </div>
        </main>
      </div>
    );
  }

  return (
    <div className="App">
      <header className="App-header">
        <h1>Lumes</h1>
        <p>Simple Lume Management</p>
      </header>
      <main className="App-content">
        <div className="lumes-section">
          <h2>Created Lumes</h2>
          <div className="lumes-grid">
            {lumes.map((lume) => (
              <LumeCard key={lume.lumeId} lume={lume} />
            ))}
          </div>
        </div>
      </main>
    </div>
  );
}

export default App;