import React, {useState, ChangeEvent, FormEvent, useEffect} from 'react';
import axios from 'axios';
import config from '../config/config'
import './SensorForm.css';


interface Location {
    latitude: number;
    longitude: number;
}
interface SensorData {
    name: string;
    tags: string[];
    description: string;
    location: Location;
}

const SensorForm: React.FC = () => {
    const initialSensorData: SensorData = {
        name: '',
        description: '',
        tags: [],
        location: {
            latitude: 0,
            longitude: 0,
        },
    };

    const [sensorData, setSensorData] = useState<SensorData>(initialSensorData);
    const [message, setMessage] = useState<string>('');

    useEffect(() => {
        if (message) {
            const timer = setTimeout(() => {
                setMessage('');
            }, 5000); // Remove message after 3 seconds

            return () => clearTimeout(timer); // Cleanup on unmount or when message changes
        }
    }, [message]);

    const handleChange = (e: ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) => {
        const { name, value } = e.target;

        if (name === 'tags') {
            const tagsArray = value.split(',').map(tag => tag.trim());
            setSensorData(prevData => ({
                ...prevData,
                tags: tagsArray,
            }));
        } else if (name === 'latitude' || name === 'longitude') {
            const locationValue = parseFloat(value);
            setSensorData(prevData => ({
                ...prevData,
                location: {
                    ...prevData.location,
                    [name]: locationValue,
                },
            }));
        } else {
            setSensorData(prevData => ({
                ...prevData,
                [name]: value,
            }));
        }
    };

    const handleReset = () => {
        setSensorData({
            name: '',
            description: '',
            tags: [],
            location: {
                latitude: 0,
                longitude: 0,
            },
        });
        setMessage('');
    };

    const handleSubmit = async (e: FormEvent<HTMLFormElement>): Promise<void> => {
        e.preventDefault();

        try {
            const response = await axios.post(`${config.apiUrl}`, sensorData);
            setMessage('Sensor data saved successfully!');
            console.log('Response:', response.data);
            setSensorData(initialSensorData); // Reset the form after successful submission
        } catch (error) {
            setMessage('An error occurred while saving sensor data.');
            console.error('Error:', error);
        }
    };

    return (
        <div className="sensor-form-container">
            <h2>Sensor Metadata Form</h2>
            <form className="" onSubmit={handleSubmit}>
                <label>
                    Name
                    <input
                        type="text"
                        name="name"
                        value={sensorData.name}
                        onChange={handleChange}
                        className="input-field"
                        required
                    />
                </label>
                <br />
                <label>
                    Description
                    <textarea
                        name="description"
                        value={sensorData.description}
                        onChange={handleChange}
                        className="input-field"
                        required
                    />
                </label>
                <br />
                <label>
                    Tags (comma-separated)
                    <input
                        type="text"
                        name="tags"
                        value={sensorData.tags.join(', ')}
                        onChange={handleChange}
                        className="input-field"
                    />
                </label>
                <br />
                <h4>Location</h4>
                <label>
                    Latitude
                    <input
                        type="number"
                        name="latitude"
                        value={sensorData.location.latitude}
                        onChange={handleChange}
                        className="input-field"
                        required
                    />
                </label>
                <br />
                <label>
                    Longitude:
                    <input
                        type="number"
                        name="longitude"
                        value={sensorData.location.longitude}
                        onChange={handleChange}
                        className="input-field"
                        required
                    />
                </label>
                <br />
                <br />
               <div style={{textAlign: "center"}}>
                   <button style={{marginRight: "10px"}} type="button" onClick={handleReset}>Reset</button>

                   <button type="submit">Save Sensor Metadata</button>
               </div>
            </form>
            {message && <p className={message.startsWith('An error') ? 'error-message' : 'success-message'}>{message}</p>}
        </div>
    );
};

export default SensorForm;
