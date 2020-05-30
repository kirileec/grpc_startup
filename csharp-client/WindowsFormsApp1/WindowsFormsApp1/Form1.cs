using Google.Protobuf;
using Grpc.Core;
using System;
using System.Collections.Generic;
using System.ComponentModel;
using System.Data;
using System.Drawing;
using System.IO;
using System.Linq;
using System.Security.Cryptography;
using System.Text;
using System.Windows.Forms;
using static Proto.FileService;

namespace WindowsFormsApp1
{
    public partial class Form1 : Form
    {
        public Form1()
        {
            InitializeComponent();
        }

        private void Form1_Load(object sender, EventArgs e)
        {

        }

        private void button1_Click(object sender, EventArgs e)
        {
            var secureChanel = new SslCredentials(File.ReadAllText("server.crt"));
            var channOptions = new List<ChannelOption>
            {
                new ChannelOption(ChannelOptions.SslTargetNameOverride,"deploy")
            };
            Channel channel = new Channel("127.0.0.1:50051", secureChanel , channOptions);

            var client = new FileServiceClient(channel);
            var req = new Proto.FSReq();
            req.DstDir = "ssdd";
            req.IfReboot = false;
            req.Name = "sadas";
            req.ProjName = "dsada";
            req.ProjType = 3;

            var file = File.ReadAllBytes("1.0-window.7z");
            SHA256Managed Sha256 = new SHA256Managed();
            byte[] bs = Sha256.ComputeHash(file);
            var hash = BitConverter.ToString(bs);
            req.Hash = hash.Replace("-","").ToLower();
            req.Filelen = file.Length;
            req.File = ByteString.CopyFrom(file);

            var reply = client.Upload(req);

            MessageBox.Show("来自" + reply.Message);

            channel.ShutdownAsync().Wait();


        }
    }
}
