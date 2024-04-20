import React from 'react';

const Algorithm = ({ label, value, checked, onChange }) => {
  return (
    <div className="flex items-center mr-4">
      <input
        type="radio"
        id={value}
        value={value}
        checked={checked}
        onChange={onChange}
        className="mr-2"
      />
      <label htmlFor={value} className="text-sm">
        {label}
      </label>
    </div>
  );
};

export default Algorithm;
