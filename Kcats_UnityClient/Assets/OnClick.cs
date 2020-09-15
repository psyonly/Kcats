using System.Collections;
using System.Collections.Generic;
using System.IO;
using System.Net.Sockets;
using UnityEngine;
using UnityEngine.UI;

public class OnClick : MonoBehaviour
{
    private bool STATUS;

    private string _ip = "127.0.0.1"; // 改为自己对外的 IP
    readonly int _port = 9874;
    string nickName = "C#";
    Socket _client;
    NetworkStream _clientStream;
    StreamWriter _clientWriter;
    StreamReader _clientReader;
    ServerCall SCManager;

    [Header("UI")]
    public InputField IF_mainText;
    public InputField IF_nickName;
    public InputField IF_socket;
    public Text chatText;
    public ScrollRect scrollRect;
    public Text T_nickName;
    public Text T_socket;
    public Text T_server;
    public Dropdown DP_clientsList;

    // datastructres
    Queue<KcatsMSG> MsgQueue;
    List<string> ClientsList;
    // Start is called before the first frame update
    void Start()
    {
        
    }

    // Update is called once per frame
    void Update()
    {
        if (!STATUS) goto RUNNING;
        if (MsgQueue.Count > 0)
        {
            KcatsMSG msg = MsgQueue.Dequeue();
            if (!msg.IsServerCall())
            {
                string text = msg.msg2Text();
                RefreshCanvas(text);
            }
            else
            {
                string[] clientsSlice = SCManager.DispatchServerCall(msg.GetMainText());
                switch (clientsSlice[0])
                {
                    case "C":
                        ClientsList.Clear();
                        DP_clientsList.ClearOptions();
                        ClientsList.Add("SYN");
                        ClientsList.Add("NIL");
                        for (int i = 1; i < clientsSlice.Length; i++)
                        {
                            ClientsList.Add(clientsSlice[i]);
                        }
                        DP_clientsList.AddOptions(ClientsList);
                        break;
                    case "E":
                        RefreshCanvas(clientsSlice[1]);
                        break;
                }
            }
        }
        RUNNING:
        if (Input.GetKeyDown(KeyCode.Return) || Input.GetKeyDown(KeyCode.KeypadEnter))
        {
            if (IF_mainText.text != "")
            {
                SendMSG();
                string addText = "\n  " + "<color=red>" + nickName + "</color>: " + IF_mainText.text;
                IF_mainText.text = "";
                IF_mainText.ActivateInputField();
                RefreshCanvas(addText);
            }
        }
    }

    public void ConfigDone()
    {
        _ip = IF_socket.text;
        nickName = IF_nickName.text;

        // 连接TCP服务器
        _client = new Socket(AddressFamily.InterNetwork, SocketType.Stream, ProtocolType.Tcp);
        _client.Connect(_ip, _port);
        _clientStream = new NetworkStream(_client);
        _clientWriter = new StreamWriter(_clientStream);
        _clientReader = new StreamReader(_clientStream);

        // 新建系统调用管理类
        SCManager = new ServerCall();

        // 显示状态
        T_server.text = _ip + ":" + _port.ToString();
        T_socket.text = _client.LocalEndPoint.ToString();
        T_nickName.text = "HostName: " + nickName;

        // 新建数据结构
        MsgQueue = new Queue<KcatsMSG>();
        ClientsList = new List<string>();

        // 异步接收消息
        ReceiveMSG();

        STATUS = true;
    }

    // 异步接收消息并处理为实例发送至管道
    public async void ReceiveMSG()
    {
        char[] buff = new char[255];
        while (true)
        {
            if(_client.Connected == false)
            {
                break;
            }
            await _clientReader.ReadAsync(buff, 0, 255);
            string msg = new string(buff);
            MsgQueue.Enqueue(DeCodeMSG(msg));
            print(msg);
        }
    }

    // 将输入框里的数据包装并发给服务器
    public void SendMSG()
    {
        //string empty = "C# Go\ncs 2 go test1\r";
        _clientWriter.WriteLine(FormMSG(DP_clientsList.captionText.text, IF_mainText.text));
        _clientWriter.Flush();
    }

    // 登出
    public void Logout()
    {
        _clientStream.Close();
        _clientWriter.Close();
        _clientReader.Close();
    }

    // 生成具有报文格式的字符串
    private string FormMSG(string recv, string text)
    {
        string str="";
        str = nickName + " " + recv + "\n" + text + "\r";
        return str;
    }

    // 根据消息字符串解码为实例
    private KcatsMSG DeCodeMSG(string msg)
    {
        string[] msgs = msg.Split('\n');
        string[] sr = msgs[0].Split(' ');
        KcatsMSG kmsg = new KcatsMSG();
        kmsg.Set(sr[0], sr[1], msgs[1]);
        return kmsg;
    }

    // 更新在线用户列表
    private void UpdateDropDown_ClientsList()
    {
        // clear & reload
        DP_clientsList.ClearOptions();
    }

    // 更新前端
    private void RefreshCanvas(string addText)
    {
        chatText.text += addText;
        Canvas.ForceUpdateCanvases();       //关键代码
        scrollRect.verticalNormalizedPosition = 0f;  //关键代码
        Canvas.ForceUpdateCanvases();   //关键代码
    }
}
