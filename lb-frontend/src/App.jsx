import React, { Fragment, useEffect, useState } from 'react'
import "axios";
import axios from 'axios';
import "./App.css";
const App = () => {
  const [subnet,setSubnet] = useState('');
  const [clientIp,setClientIp] = useState('');
  const [lbIp,setLbIp] = useState('');
  const [serverIp,setServerIp] = useState('');
  const [getsn,setGetsn]= useState('');
  const [info,setInfo] = useState({});
  // for /api/create_subnet
  const handleSubnet = async (e) =>{
    e.preventDefault();
    const data = {
      subnet,
    };
    const headers = {
      'Content-Type':'application/json',
    }
    const url = "http://localhost:8080/api/create_subnet";
    const res = await axios.post(url,data,{headers});
    alert(res.data.message);
    window.location.reload();
  };
  // for /api/launch_client
  const handleClientIp = async (e)=>{
    e.preventDefault();
    const data = {
      ip:clientIp,
    };
    const headers = {
      'Content-Type':'application/json',
    }
    const url = "http://localhost:8080/api/launch_client";
    const res = await axios.post(url,data,{headers});
    alert(res.data.message);
    window.location.reload();
  };
  // for /api/launch_lb
  const handleLbIp = async (e) => {
    e.preventDefault();
    const data = {
      ip:lbIp,
    };
    const headers = {
      'Content-Type':'application/json',
    }
    const url = "http://localhost:8080/api/launch_lb";
    const res = await axios.post(url,data,{headers});
    alert(res.data.message);
    window.location.reload();
  };
  // for /api/launch_server
  const handleServerIp = async (e) => {
    e.preventDefault();
    const data = {
      ip:serverIp,
    };
    const headers = {
      'Content-Type':'application/json',
    }
    const url = "http://localhost:8080/api/launch_server";
    const res = await axios.post(url,data,{headers});
    alert(res.data.message);
    window.location.reload();
  };
  // for /api/attach_xdp
  const handleAttachXdp = async (e)=>{
    e.preventDefault();
    const url = "http://localhost:8080/api/attach_xdp";
    const res = await axios.get(url);
    alert(res.data.message);
  }

  // function to unparse the IP address
 const unparseIP = (parsedIP) => {
      // Convert the uint32 IP to an array of 4 bytes
      const bytes = [
        (parsedIP >> 0) & 0xFF,
        (parsedIP >> 8) & 0xFF,
        (parsedIP >> 16) & 0xFF,
        (parsedIP >> 24 )& 0xFF
      ];
      // Join the bytes with dots to form the IP address string
      return bytes.join('.');
    }
  // function to unparse the MAC address
  const unparseMAC = (parsedMAC) => {
    // Convert the array of bytes to a MAC address string
    return parsedMAC.map(byte => byte.toString(16).padStart(2, '0')).join(':');
  }
  
  useEffect(()=>{
    // get the map data on window load
    // for /api/get_data
    const fetchInfo = async ()=> {
      const url = "http://localhost:8080/api/get_data";
      const res = await axios.get(url);
      setInfo(res.data.data);
      setGetsn(res.data.subnet);
    
    };
    fetchInfo();
  },[])
  return (
    <Fragment>
      <div className='container'>
          <div className='left-panel'>
            <div className='subnet-input'>
              <label>SUBNET </label>
              <input type="text" placeholder={"192.168.0.0/16"} value={subnet} onChange={(e)=>{setSubnet(e.target.value)}}></input>
              <button onClick={(e)=> handleSubnet(e)}>Create Subnet</button>
            </div>
            <div className='clientip-input'>
              <label>CLIENT IP (max entry 1)</label>
              <input type="text" placeholder={"192.168.0.3"} value={clientIp} onChange={(e)=>{setClientIp(e.target.value)}}></input>
              <button onClick={(e)=> handleClientIp(e)}>Launch Client</button>
            </div>
            <div className='lb-input'>
              <label>LOAD-BALANCER IP (max entry 1)</label>
              <input type="text" placeholder={"192.168.0.6"} value={lbIp} onChange={(e)=>{setLbIp(e.target.value)}}></input>
              <button onClick={(e)=> handleLbIp(e)}>Launch LoadBalancer</button>
            </div>
            <div className='server-input'>
              <label>SERVER IP (max entry 128)</label>
              <input type="text" value={serverIp} placeholder={"192.168.0.7"} onChange={(e)=>{setServerIp(e.target.value)}}></input>
              <button onClick={(e)=> handleServerIp(e)}>Launch Server</button>
            </div>
            <div className='attach-xdp'>
              <button onClick={(e)=> handleAttachXdp(e)}>Attach XDP</button>
            </div>
          </div>
          <div className='right-panel'>
              <div className="subnet-show">
                {
                  getsn? (<h2>SUBNET - {getsn}</h2>) : (<h2>Subnet Will Appear Here</h2>)
                  }
              </div>
              <div className="show-table">
                    <table>
                      <thead>
                        <tr>
                          <th>IP Address</th>
                          <th>MAC Address</th>
                        </tr>
                      </thead>
                      <tbody>
                       {Object.entries(info).map(([key, value]) => (
                          <tr key={key}>
                            <td>{unparseIP(value.Ip)}</td>
                            <td>{unparseMAC(value.Mac.Addr)}</td>
                    
                          </tr>))}
                      </tbody>
                    </table>
              </div>
          </div>
      </div>
    </Fragment>
  )
}

export default App