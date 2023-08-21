import React from 'react';
import { render, fireEvent } from '@testing-library/react';
import '@testing-library/jest-dom/extend-expect';
import SensorForm from './SensorForm';

test('renders SensorForm with form fields and button', () => {
    // @ts-ignore
    const { getByText, getByLabelText, getByTestId } = render(<SensorForm />);

    // Check if the form fields are rendered
    const sensorNameInput = getByLabelText(/Name/i);
    const tagsInput = getByLabelText(/Tags/i);
    const descriptionInput = getByLabelText(/Description/i);
    const latitudeInput = getByLabelText(/Latitude/i);
    const longitudeInput = getByLabelText(/Longitude/i);

    expect(sensorNameInput).toBeInTheDocument();
    expect(tagsInput).toBeInTheDocument();
    expect(descriptionInput).toBeInTheDocument();
    expect(latitudeInput).toBeInTheDocument();
    expect(longitudeInput).toBeInTheDocument();

    // Check if the button is rendered
    const saveButton = getByText(/Save Sensor Metadata/i);
    expect(saveButton).toBeInTheDocument();

    // Simulate user input and button click
    fireEvent.change(sensorNameInput, { target: { value: 'Test Sensor' } });
    fireEvent.change(tagsInput, { target: { value: 'tag1, tag2' } });
    fireEvent.change(descriptionInput, { target: { value: 'Test description' } });
    fireEvent.change(latitudeInput, { target: { value: 42.123 } });
    fireEvent.change(longitudeInput, { target: { value: -71.456 } });

    fireEvent.click(saveButton);

});


