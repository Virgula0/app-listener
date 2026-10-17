# Graph Report - app-listener  (2026-10-17)

## Corpus Check
- 457 files · ~1,639,221 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 13 file(s) not represented in the graph (top: (none) 8, .service 2, .pub 1)

## Summary
- 5931 nodes · 25212 edges · 219 communities (184 shown, 35 thin omitted)
- Extraction: 91% EXTRACTED · 9% INFERRED · 0% AMBIGUOUS · INFERRED: 2313 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `bd161150`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- daemon/daemon.go
- testing.T
- IntegrationSuite
- SampleContent
- daemonconfig_test.go
- networkguard.bpf.c
- monitor.bpf.c
- gui/session.go
- guardUnitTest
- mkdirs
- codeedit.go
- auditGroupCoverage
- protectedPage
- github.com/charmbracelet/bubbletea.Cmd
- .launch
- sync.Mutex
- IntegrationSuite
- guard_trust.bpf.c
- IntegrationSuite
- fileEditModel
- networkmonitor.bpf.c
- usecase/daemon_test.go
- Resource
- auth.go
- runNetworkGuard
- guard.bpf.c
- SelectOrphans
- vaultfs.go
- icon
- selection.go
- UI
- serve.go
- IntegrationSuite
- github.com/charmbracelet/bubbles/textarea.Model
- picker
- github.com/spf13/cobra.Command
- time.Duration
- backend
- controlServer
- fileedit_model_test.go
- os.File
- IntegrationSuite
- absPath
- monitorUnitTest
- Config
- update/update.go
- Ledger
- stdio
- editView
- runDaemon
- LibraryClosure
- GenerateConf
- Client
- watchGroup
- guardModel
- TrustGuard
- systemd.go
- update_test.go
- GuardInodeKey
- Vault
- buildOneGuard
- pill
- exe_supersede.h
- runEditProtected
- time.Time
- bpfstats/main.go
- wm.go
- runMonitor
- fakeGuardRepo
- FileTarget
- fcntl
- GuardRepository
- Guard
- watchSet
- btrfsLayoutFrom
- changelog_test.go
- tempGrantSetup
- RefreshSection
- NetEventType
- process_vm_readv.c
- GuardEvent
- DaemonEvent
- planUpdaters
- User
- filevault_test.go
- buildTrustedSet
- find.go
- daemon_multicall_test.go
- liveSession
- daemon mode (fscrypt + whitelist lifecycle)
- FORWARD temporary -w/-b whitelist grants
- runUninstall
- inodeChain
- netinject.c
- NetGuard
- check-compatibility.sh
- BPF verifier 1M-insn complexity budget
- Inode-based (dev:ino) executable identity
- CandidateDir
- policyEditor
- CopyTreeWithProgress
- Well-known bypass classes
- NewDumpHook
- jit_provenance.c
- inspector_probe.c
- ptrace_race.c
- downloadFile
- Build & sign release assets (reusable workflow)
- installPage
- KernelDev
- mountTable
- NetworkGuardRepository
- app-listener Terminal Demo Recording
- IntegrationSuite
- shQuote
- changelog.go
- IntegrationSuite
- IntegrationSuite
- Regenerate eBPF bindings in pinned builder image
- serveWebSocket
- IntegrationSuite
- WebSocket /ws client
- launch_probe.c
- install.sh
- mc_tag.h
- .FilterExistingWhitelist
- make test-integration (rootful Docker suite)
- fakeDaemon
- supersede_probe.c
- newTestVetter
- trace-app-libs.sh
- IntegrationSuite
- configedit_test.go
- Bun TMPDIR redirect with reserved .bun-* name
- Daemon self-protection (selfguards.go)
- github.com/Virgula0/app-listener
- network-monitor mode
- uninstall subcommand
- TempBinary
- NewEventFanout
- SanitizeText
- hashFileShared
- catalogRefresher
- net_tester/main.go
- edit-protected
- Monitor
- hashCache
- daemon: permanent protection with encryption at rest
- ConfinedEntry
- EventType
- .Layout
- RunWith
- newTestTarget
- Installation
- backend_test.go
- graphify-refresh.sh
- writeWithin
- golang.org/x/sys/unix.Stat_t
- main_test.go
- pinstate.go
- feed.go
- shell
- integrationtests/networkguard_test.go
- WithBottomBar
- IntegrationSuite
- lockRootRecovering
- IntegrationSuite
- Limitations
- Troubleshooting
- IntegrationSuite
- Vault
- Daemon configuration
- Active
- batchInodeRows
- CheckBPFLSM
- Compatibility
- How it works
- ResolveTempBinary
- editControlSession
- bleedLayout
- Parse
- Binary trust
- guard: block file access by program
- network-guard: block network access by program
- lockAndDeprovision
- Development
- app-listener documentation
- Self-updating apps and the catalog
- network-monitor: see what a program talks to
- mountsUnitTest
- github.com/stretchr/testify/suite.Suite
- bufio.Reader
- monitor: see what touches a path
- kernel-prof.sh
- pathCache
- .eventTypeAllowed
- Execute
- decryptStep
- .populateWatchInodes
- WithSuspectsFrom
- copyRegularAt
- Desktop GUI
- ParseEventsFlag
- runLockdown
- parseResize
- checkRunningBinaryMatchesInstalled
- TestSuperblockDevMatchesStatOffBtrfs
- TestCleanRel
- bytes.Buffer
- golang.org/x/sys/unix.Timespec
- MulticallError

## God Nodes (most connected - your core abstractions)
1. `Guard` - 117 edges
2. `Config` - 96 edges
3. `GuardInodeKey` - 96 edges
4. `Load()` - 90 edges
5. `IntegrationSuite` - 82 edges
6. `absPath()` - 81 edges
7. `shell` - 81 edges
8. `fileEditModel` - 77 edges
9. `writeConfig()` - 69 edges
10. `IntegrationSuite` - 68 edges

## Surprising Connections (you probably didn't know these)
- `Editor keys` --references--> `Delete()`  [INFERRED]
  docs/edit-protected.md → internal/backups/backups.go
- `Well-known bypass classes` --references--> `ptrace_race POC`  [INFERRED]
  CLAUDE.md → integrationtests/exploits/README.md
- `lib_probe reserved library POC` --semantically_similar_to--> `Bun TMPDIR redirect with reserved .bun-* name`  [INFERRED] [semantically similar]
  integrationtests/exploits/README.md → README.md
- `swap_benign / swap_reader in-place binary swap POC` --semantically_similar_to--> `Binary ledger (/etc/app-listener/binaries.db)`  [INFERRED] [semantically similar]
  integrationtests/exploits/README.md → README.md
- `resolveAssets()` --calls--> `asset()`  [INFERRED]
  cmd/functions/update/update.go → internal/gui/icons.go

## Import Cycles
- None detected.

## Communities (219 total, 35 thin omitted)

### Community 0 - "daemon/daemon.go"
Cohesion: 0.09
Nodes (14): tempJournalRow, User, launchRule, trustEvent, textAnswer, btrfsIoctlFsInfoArgs, daemonStatsMsg, guardStatsMsg (+6 more)

### Community 1 - "testing.T"
Cohesion: 0.02
Nodes (145): TestDaemonSelfBaselineEventsNeverWrites(), TestGroupUnlockSelfEventsIsBaselinePlusMknod(), TestDiffCatalogFlagRegistered(), TestMovePeer(), TestParseAllowLibNotADirective(), TestParseLibBinaryNotADirective(), TestParseLibDirNotADirective(), helper() (+137 more)

### Community 2 - "IntegrationSuite"
Cohesion: 0.09
Nodes (9): configBinaryLines(), IntegrationSuite, parseDaemonEvents(), probeExeAs(), procGateMatrix(), procReadProbe(), rawExec(), watchPathsInConfig() (+1 more)

### Community 3 - "SampleContent"
Cohesion: 0.11
Nodes (21): TestDaemonUnitWithMetadataOutput(), detectSSHAgentUnits(), isInstallerSSHAgentUnit(), removeSSHAgentUnit(), TestIsInstallerSSHAgentUnit(), TestRemoveSSHAgentUnit(), TestDaemonUnitCmdPrintsBundledUnit(), SampleContent() (+13 more)

### Community 4 - "daemonconfig_test.go"
Cohesion: 0.08
Nodes (73): Load(), ResolvePendingPaths(), hasBinary(), hasPendingBinary(), TestElectronAppsBlockParsed(), TestElectronAppsBlockRejectsMalformedLines(), TestEncryptionGroupsSkipsLibDirs(), TestInspectorsBlockBounded() (+65 more)

### Community 5 - "networkguard.bpf.c"
Cohesion: 0.10
Nodes (48): admit_refusal(), check_watched(), code_vouched(), emit_event(), emit_gate(), env_loader(), env_step(), exe_inode_of() (+40 more)

### Community 6 - "monitor.bpf.c"
Cohesion: 0.15
Nodes (32): emit_event(), emit_event_kern(), fill_identity(), fill_path(), read_path(), trace_do_faccessat(), trace_do_sendfile(), trace_do_splice() (+24 more)

### Community 7 - "gui/session.go"
Cohesion: 0.05
Nodes (67): installConfig(), installFile(), installFileAs(), installUnit(), readAtNoFollow(), upsertFile(), upsertUnitAt(), writeAtNoFollow() (+59 more)

### Community 8 - "guardUnitTest"
Cohesion: 0.05
Nodes (23): addErrKind, ReplacementCheck, BinariesSummary(), classifyAddErr(), guardModeKey(), addRecorder(), copySelf(), guardUnitTest (+15 more)

### Community 9 - "mkdirs"
Cohesion: 0.14
Nodes (35): newWatchPattern(), newFakeInotify(), TestWatchSetBinaryWriteTriggersNotMkdir(), TestWatchSetFloodIsCapped(), TestWatchSetLibDirAppearanceTriggers(), TestWatchSetMissingRootWatchesDeepestAncestor(), TestWatchSetOverflowRescansAndRebuilds(), TestWatchSetSymlinkCreateTriggers() (+27 more)

### Community 10 - "codeedit.go"
Cohesion: 0.10
Nodes (39): class, classOf(), Edit(), leadingIndent(), Navigate(), New(), onlySpaces(), SetText() (+31 more)

### Community 11 - "auditGroupCoverage"
Cohesion: 0.22
Nodes (11): auditEntry(), auditFindingsWithConfig(), auditGroupCoverage(), auditTree(), coveredByWatch(), findResource(), groupWatchPaths(), relOrPath() (+3 more)

### Community 12 - "protectedPage"
Cohesion: 0.05
Nodes (25): backendLink, configView, forwardForm, forwardView, installUI, multiAnswer, protectedPage, Event (+17 more)

### Community 13 - "github.com/charmbracelet/bubbletea.Cmd"
Cohesion: 0.06
Nodes (16): snapNetGuard, diffModel, editorModel, NetGuardEvent, protoString(), renderBullets(), TestNoticeViewFitsAndPlacesOKLast(), TestRenderBulletsVerbatimAndIndented() (+8 more)

### Community 14 - ".launch"
Cohesion: 0.27
Nodes (3): IntegrationSuite, launchCase, TestClientVerifyRunsBeforeAnythingIsSent()

### Community 15 - "sync.Mutex"
Cohesion: 0.09
Nodes (11): netGuardPill(), newActivity(), TestRingKeepsNewest(), clonePolicy(), NetExtra, NetPolicy, Network, Counts (+3 more)

### Community 16 - "IntegrationSuite"
Cohesion: 0.10
Nodes (9): attrOpCase, eventsForPath(), IntegrationSuite, guardDeltaEvents(), guardEventTypes(), le64HexKey(), parseGuardEvents(), guardEvent (+1 more)

### Community 17 - "guard_trust.bpf.c"
Cohesion: 0.09
Nodes (57): caller_updates(), current_exe_flags(), current_exe_inode(), current_exe_key(), current_writer_bits(), dentry_is_dir(), dentry_parent(), deny_if_protected() (+49 more)

### Community 18 - "IntegrationSuite"
Cohesion: 0.16
Nodes (7): netMonitorEvent, IntegrationSuite, fmtEvents(), IntegrationSuite, netMonitorDeltaEvents(), netMonitorHasTypes(), parseNetMonitorEvents()

### Community 19 - "fileEditModel"
Cohesion: 0.08
Nodes (14): clipLabel(), humanSize(), listEntries(), clamp(), paneTitle(), sortNodes(), maxLineWidth(), TestFileEditPaneGeometry() (+6 more)

### Community 20 - "networkmonitor.bpf.c"
Cohesion: 0.23
Nodes (16): BPF_KRETPROBE(), emit_event(), get_netns(), get_socket_proto(), is_watched_binary(), read_inet_addr(), trace_accept(), trace_accept4() (+8 more)

### Community 21 - "usecase/daemon_test.go"
Cohesion: 0.14
Nodes (49): NewDaemonUseCase(), partitionEncryptionRoots(), newFakeVault(), resource(), startDaemon(), TestDaemonUseCaseConstructorMismatch(), TestDaemonUseCaseEventsMerged(), TestDaemonUseCaseGrantEditAccess() (+41 more)

### Community 22 - "Resource"
Cohesion: 0.13
Nodes (8): Resource, concurrencyLimit(), encryptionRootOf(), DaemonUseCase, uniqueEncryptionRoots(), unlockRoots(), vaultOpForGuard(), vaultOpOnRoot()

### Community 23 - "auth.go"
Cohesion: 0.16
Nodes (23): Hash(), HashFileExists(), LoadHashFile(), OriginOf(), parseHash(), RemoveHashFile(), hashFileInode(), TestHashFileLifecycle() (+15 more)

### Community 24 - "runNetworkGuard"
Cohesion: 0.09
Nodes (30): CheckEBPF(), ParseNetEventsFlag(), prepareDaemonStart(), computeGuardBinaries(), confirmUnsafe(), discoverInfra(), resolveGuardMode(), runNetworkGuard() (+22 more)

### Community 25 - "guard.bpf.c"
Cohesion: 0.09
Nodes (72): inode_dev(), add_inode_to_guard(), btrfs_copy_gate(), check_and_emit_args(), chmod_only_drops_write(), code_suspect(), count_degrade(), current_is_inspector() (+64 more)

### Community 26 - "SelectOrphans"
Cohesion: 0.23
Nodes (16): appProtectorSet(), CleanOrphans(), isAppProtector(), listPolicies(), listProtectors(), orphanedPolicies(), orphanedProtectors(), SelectOrphans() (+8 more)

### Community 27 - "vaultfs.go"
Cohesion: 0.09
Nodes (24): DefaultNewFileMode(), Vault, IsBinaryContent(), IsOctalMode(), Open(), ParseMode(), inodeOf(), mustRead() (+16 more)

### Community 28 - "icon"
Cohesion: 0.08
Nodes (32): binaryList, checkGrid, filesPage, handButton, homePage, panel, statusTile, newBinaryList() (+24 more)

### Community 29 - "selection.go"
Cohesion: 0.14
Nodes (16): pos, lines(), cellToRune(), clipboardText(), cursorPos(), Selection, offset(), setClipboard() (+8 more)

### Community 30 - "UI"
Cohesion: 0.13
Nodes (9): MultiSelectKeymap(), Confirm(), Choice, Field, Question, UI, Input(), MultiSelect() (+1 more)

### Community 31 - "serve.go"
Cohesion: 0.16
Nodes (18): basicAuth(), isInteractiveTerminal(), newServeListener(), newServeServer(), relayServeSignals(), runServe(), securityHeaders(), Serve() (+10 more)

### Community 33 - "github.com/charmbracelet/bubbles/textarea.Model"
Cohesion: 0.11
Nodes (18): lineStyles, viewportState, column(), dedent(), moveTo(), repeatKey(), reposition(), Finder (+10 more)

### Community 34 - "picker"
Cohesion: 0.06
Nodes (37): keyEntry, listOpts, pickEntry, picker, pickMode, treeSearch, openPicker(), pickerTree() (+29 more)

### Community 35 - "github.com/spf13/cobra.Command"
Cohesion: 0.12
Nodes (23): AddServeFlags(), credentialFlagState(), isLoopbackHost(), ParseServeFlags(), requestedBoolFlag(), validateEnabledServe(), validateServeAddress(), TestParseServeFlagsRequiresCredentials() (+15 more)

### Community 36 - "time.Duration"
Cohesion: 0.15
Nodes (8): procStats, statsSample, IntegrationSuite, pooledContainer, cpuTime(), humanBytes(), rssBytes(), Ingested()

### Community 37 - "backend"
Cohesion: 0.09
Nodes (18): TestValidatePassword(), ValidatePassword(), checkSecret(), entryOf(), TestCheckSecretRefusesLineBreaks(), revertBunLaunchers(), File, backend (+10 more)

### Community 38 - "controlServer"
Cohesion: 0.24
Nodes (5): readConfigPut(), startControlServer(), TestRecordFailureLockout(), configEditor, controlServer

### Community 39 - "fileedit_model_test.go"
Cohesion: 0.21
Nodes (23): newFileEditModel(), inodeOf(), selectNamed(), TestChmodSuidEndToEnd(), TestFileChmod(), TestFileChownFlow(), TestFileCreateAndDeleteFlow(), TestFileEditModelSingleFileCreateRefused() (+15 more)

### Community 40 - "os.File"
Cohesion: 0.05
Nodes (59): checkFreshDir(), fileID(), openBunTmp(), openDirNoFollow(), renewDir(), plantBun(), TestOpenBunTmp_KeepsVettedDir(), TestOpenBunTmp_ReplacesUnvettedDir() (+51 more)

### Community 41 - "IntegrationSuite"
Cohesion: 0.20
Nodes (3): exploitTest, newEventTypes(), IntegrationSuite

### Community 42 - "absPath"
Cohesion: 0.08
Nodes (4): IntegrationSuite, inNosuidNamespace(), IntegrationSuite, absPath()

### Community 44 - "Config"
Cohesion: 0.11
Nodes (24): newBinaryVetter(), openBinaryVetter(), newCatalogRefresher(), makeReloadHandler(), noteReloadedSettings(), startCatalogRefresh(), resolveInspectors(), TestResolveInspectorsKeepsOnlyAdmitted() (+16 more)

### Community 45 - "update/update.go"
Cohesion: 0.16
Nodes (23): applyUpdate(), assetsFor(), channelOf(), checkUpdatePreconditions(), confirmUnit(), deploy(), downloadAndVerify(), downloadReleaseFiles() (+15 more)

### Community 46 - "Ledger"
Cohesion: 0.11
Nodes (11): candidates(), confirmable(), confirmAll(), describe(), hashNow(), run(), fyneLogFilter, Ledger (+3 more)

### Community 47 - "stdio"
Cohesion: 0.07
Nodes (4): put(), wait_then_read(), hold_forever(), main()

### Community 48 - "editView"
Cohesion: 0.12
Nodes (8): codeEntry, editView, Entry, selectable(), fileSize(), joinRel(), newCodeEntry(), relOf()

### Community 49 - "runDaemon"
Cohesion: 0.12
Nodes (13): newControlManager(), startControlManager(), awaitStartupOrSignal(), catchLifecycleSignals(), runDaemon(), startGuardedDaemonAbortable(), startPprof(), writePidFile() (+5 more)

### Community 50 - "LibraryClosure"
Cohesion: 0.10
Nodes (19): ldConfParser, defaultLibDirs(), elfInterp(), fileExists(), isScript(), ldConfInclude(), ldSoConfDirs(), LdSoPreloadPaths() (+11 more)

### Community 51 - "GenerateConf"
Cohesion: 0.05
Nodes (64): nextConfig(), appendSectionsAndEdit(), TestDiffMergePreservesExistingAndParses(), collectFilesystemPrereqs(), askEncryption(), editConfig(), libraryBlocksFromCandidates(), runConfigEditor() (+56 more)

### Community 52 - "Client"
Cohesion: 0.07
Nodes (21): Client, stagedBinary, childEnv(), launch(), Start(), verifiedBinary(), TestChildEnvIsAllowlisted(), TestHelperBackend() (+13 more)

### Community 53 - "watchGroup"
Cohesion: 0.10
Nodes (36): watchGroup, addResource(), applyAllowLib(), applyConfigLine(), applyDirective(), applyElectronApp(), applyInspector(), applyLibBinary() (+28 more)

### Community 54 - "guardModel"
Cohesion: 0.06
Nodes (33): Stats, parseKBValue(), parseStat(), parseStatus(), Read(), TestReadConcurrentRace(), TestStatsFieldsNonZero(), TestStatsValuesMonotonic() (+25 more)

### Community 55 - "TrustGuard"
Cohesion: 0.05
Nodes (20): exeRow, tempGrant, tempMaskOp, tempRow, trustHook, trustRows, deleteInoKeys(), deleteResKeys() (+12 more)

### Community 56 - "systemd.go"
Cohesion: 0.12
Nodes (33): reloadDaemonForPasswordChange(), buildBinaryIfNeeded(), installBinaryOnly(), mustCwd(), revertSystemFiles(), DaemonUnitInstalled(), DeployInstalledBinary(), DeployInstalledBinaryNoRestart() (+25 more)

### Community 57 - "update_test.go"
Cohesion: 0.10
Nodes (14): parseChecksum(), parsePublicKey(), rsaPublicKeyPEM(), signedChecksum(), TestAssetsFor(), TestFetchReleaseByTag(), testKeyPair(), TestParseChecksum() (+6 more)

### Community 58 - "GuardInodeKey"
Cohesion: 0.03
Nodes (51): engine, refusedReplacements, resInfo, taintOwner, TestNearestRootClaims(), TestTieRank(), nextLabel(), TestDeliverLabelsCrossResourceProcessGates() (+43 more)

### Community 59 - "Vault"
Cohesion: 0.18
Nodes (13): isFileVaultCiphertext(), isRegularFileTarget(), Vault, hasEncryptionPolicy(), isLockedRegularFileErr(), isNotEncryptedErr(), newBoundedKeyFn(), readKey() (+5 more)

### Community 60 - "buildOneGuard"
Cohesion: 0.06
Nodes (54): backingDeviceUnion(), buildConcurrency(), buildGuards(), buildOneGuard(), eventFilterOptions(), newPinGeneration(), reloadOnce(), reloadSlotsNeeded() (+46 more)

### Community 61 - "pill"
Cohesion: 0.03
Nodes (17): appTheme, binaryRow, clickPanel, iconTile, iconTileRenderer, navItem, navItemRenderer, panelRenderer (+9 more)

### Community 62 - "exe_supersede.h"
Cohesion: 0.14
Nodes (20): exe_is_superseded(), exe_last_link_key(), exe_note_freed(), exe_note_supersede(), exe_refused(), exe_seq_now(), exe_stamp_exec(), exe_stamp_fork() (+12 more)

### Community 63 - "runEditProtected"
Cohesion: 0.09
Nodes (24): auditAfterEdit(), auditFindings(), reserveStdout(), runBackend(), LiveModeAvailable(), runEditProtected(), TestPickEncryptedDirSingleSkipsPicker(), pingInterval() (+16 more)

### Community 64 - "time.Time"
Cohesion: 0.07
Nodes (14): sweepID, sweepMemo, logHook, logLine, runHost, runView, stageItem, stageState (+6 more)

### Community 65 - "bpfstats/main.go"
Cohesion: 0.15
Nodes (21): guardSpec(), TestGuardSpecKeepsIoctlCompatOnlyWhereKernelHasIt(), verifyCollection(), VerifyLoad(), applyInodeCapacity(), TestApplyInodeCapacity(), InertProgram(), KernelHasFunc() (+13 more)

### Community 66 - "wm.go"
Cohesion: 0.11
Nodes (25): detectWM(), reachable(), centered(), Detect(), FitSize(), Session, hyprlandSocket(), i3RuntimeSocket() (+17 more)

### Community 67 - "runMonitor"
Cohesion: 0.11
Nodes (18): BuildEBPFTargets(), RawTarget, isSubDir(), MakeDisplayPaths(), pathWithinDir(), ResolveTargets(), validatePair(), ValidateTargets() (+10 more)

### Community 69 - "FileTarget"
Cohesion: 0.12
Nodes (13): Target, specFromForm(), guardPill(), specSummary(), EventTypes(), Files, FileSpec, FileTarget (+5 more)

### Community 70 - "fcntl"
Cohesion: 0.09
Nodes (9): main(), report(), denied(), main(), dump_via_debugfs(), main(), probe_open(), denied() (+1 more)

### Community 72 - "Guard"
Cohesion: 0.05
Nodes (31): deferredBinary, rootHandle, snapNet, BinaryEntry, openRoot(), TestCheckBinaryEventsReadOnlyAcceptsUnrestricted(), checkBtrfsKeys(), SharedPinDegraded() (+23 more)

### Community 73 - "watchSet"
Cohesion: 0.12
Nodes (17): refreshAdmits(), TestRefreshAdmitsHomeMatchesAndRootPlacedNamesOnly(), addSystemBinary(), catalogPatterns(), catalogWatchPlan(), inHome(), isSymlink(), newWatchSet() (+9 more)

### Community 74 - "btrfsLayoutFrom"
Cohesion: 0.14
Nodes (12): BtrfsLayout, fakeTypes, typeSource, ensureBtrfsLayout(), btrfsLayoutFrom(), findMember(), memberOffset(), ResolveBtrfsLayout() (+4 more)

### Community 75 - "changelog_test.go"
Cohesion: 0.25
Nodes (10): newChangelogModel(), rel(), tags(), TestChangelogModelOtherKeysDoNotQuit(), TestChangelogModelQuitKeys(), TestChangelogModelResize(), TestChangelogText(), TestCumulativeChangelog() (+2 more)

### Community 76 - "tempGrantSetup"
Cohesion: 0.29
Nodes (8): journalTo(), tempBinary(), tempGrantSetup(), TestGrantTemporaryAccessJournalsBeforeApplyAndClearsAfterRevoke(), TestGrantTemporaryAccessRefusals(), TestGrantTemporaryAccessRevokesOnInPlaceRewrite(), TestGrantTemporaryAccessRollsBackOnApplyFailure(), tempTrace

### Community 77 - "RefreshSection"
Cohesion: 0.11
Nodes (25): logRefreshChange(), sameConfig(), discordRefreshFixture(), TestRefreshOnceAppliesAndReloads(), TestRefreshOnceKeepsOperatorLinesAndSkipsNoOps(), TestRefreshOnceRestoresTheFileWhenTheReloadFails(), TestRefreshOnceSkipsAConfigChangedSinceLoad(), TestLogRefreshChange() (+17 more)

### Community 78 - "NetEventType"
Cohesion: 0.06
Nodes (18): snapNetMonitor, NetBpfEvent, NetEvent, NetEventType, newNetEventChecks(), FormatAddr(), NetEventTypes(), Ntohs() (+10 more)

### Community 79 - "process_vm_readv.c"
Cohesion: 0.18
Nodes (6): dump_memory(), find_all_rw_regions(), find_heap(), find_victim_by_fd(), main(), spawn_less_pty()

### Community 80 - "GuardEvent"
Cohesion: 0.05
Nodes (27): newTestSession(), TestRunSessionIdleTimeoutResetsOnPing(), TestRunSessionStreamsEventsAsActivity(), holdHeadless(), HeadlessLine(), runGuardHeadless(), bpfGuardEvent, snapGuard (+19 more)

### Community 81 - "DaemonEvent"
Cohesion: 0.08
Nodes (37): ServeConfig, UIDResolver, NewUIDResolver(), drainEphemeral(), newDaemonModel(), runDaemonUI(), runHeadless(), runServedTUI() (+29 more)

### Community 82 - "planUpdaters"
Cohesion: 0.16
Nodes (16): hardLinked(), libBinaries(), mayUpdate(), planUpdaters(), resourceBinaries(), resourceLinks(), rules(), TestGeneralToolsToWarn_HardLinkIsNoUpdaterButNoWarning() (+8 more)

### Community 83 - "User"
Cohesion: 0.06
Nodes (54): askBunTmpdirUsers(), bunEntryConfigured(), bunLaunchersFor(), gatherPerUserSetup(), setupBunTmpdir(), applyDiffAdditions(), collectDiffAdditions(), configPaths() (+46 more)

### Community 84 - "filevault_test.go"
Cohesion: 0.10
Nodes (46): decryptDirectories(), TestDecryptDirectoriesEmpty(), classifyRegularFileTarget(), clearRecovery(), deriveFileVaultSubkey(), EnsureRecoverySidecarPlaceholder(), fileVaultRecoverPath(), looksLikeFileVaultRecord() (+38 more)

### Community 85 - "buildTrustedSet"
Cohesion: 0.07
Nodes (30): inspectorPaths(), TestInspectorPathsOnlyRootPlaced(), buildTrustedSet(), TestBuildTrustedSetExcludesLibraryClosure(), closureRejections(), inGuardedTree(), keys(), TestWarnUntrustedLibs_OnlyWhatTheKernelRefuses() (+22 more)

### Community 86 - "find.go"
Cohesion: 0.18
Nodes (18): findButton, span, findAll(), lowerRunes(), NewFinder(), newFindInput(), overlay(), runesEqual() (+10 more)

### Community 87 - "daemon_multicall_test.go"
Cohesion: 0.17
Nodes (10): decodeMulticall(), IntegrationSuite, multicallBypasses(), multicallConfig(), uutilsStrays(), vettedGroupIndex(), tailLast(), mcBypass (+2 more)

### Community 88 - "liveSession"
Cohesion: 0.17
Nodes (11): parseErr(), TestLiveSessionConfigAndPut(), TestLiveSessionWatchAndPinger(), chooseResources(), displayEntries(), eventsLabel(), forwardBinaries(), holdForward() (+3 more)

### Community 89 - "daemon mode (fscrypt + whitelist lifecycle)"
Cohesion: 0.18
Nodes (10): Daemon catalog refresh (catalogrefresh.go / catalogwatch.go), exe_supersede.h superseded-key tracking, glob_plant reserved glob name POC, launch_probe risky-launch POC, preload_wait code-suspect POC, supersede_probe reused-inode POC, Daemon boot ordering (Type=notify, Before=user sessions), daemon mode (fscrypt + whitelist lifecycle) (+2 more)

### Community 90 - "FORWARD temporary -w/-b whitelist grants"
Cohesion: 0.21
Nodes (9): CONFIG post-AUTH daemon.conf editing (configedit.go), FORWARD temporary -w/-b whitelist grants, Daemon self-key (GUARD_ALLOW_ROOT) bypass, swap_benign / swap_reader in-place binary swap POC, edit-protected --edit-config, edit-protected subcommand, edit-protected live mode, edit-protected --forward (+1 more)

### Community 91 - "runUninstall"
Cohesion: 0.17
Nodes (21): deletePostBackups(), restoreBackups(), offerBackupCleanup(), runUninstall(), removeKeyAndEmptyDir(), removeMasterKey(), removeSSHAgentEnv(), revertSSHAgents() (+13 more)

### Community 92 - "inodeChain"
Cohesion: 0.21
Nodes (16): fdStat, inodeChain(), isProcFD(), mountinfoShowsSbRoot(), mountShowsSbRoot(), physicalParent(), statFD(), keyOf() (+8 more)

### Community 93 - "netinject.c"
Cohesion: 0.17
Nodes (10): attach_child(), await_exec(), do_connect(), errname(), main(), sleep_ms(), traceme(), do_connect() (+2 more)

### Community 95 - "NetGuard"
Cohesion: 0.09
Nodes (12): foreignPinnedLink(), reasonText(), mcNameKey(), statInodeKey(), judgePasses(), parseProcs(), TestJudgePasses(), TestParseProcs() (+4 more)

### Community 96 - "check-compatibility.sh"
Cohesion: 0.38
Nodes (11): fail(), load_config(), lsm_bpf_instructions(), note(), pass(), read_config(), require_config(), section() (+3 more)

### Community 97 - "BPF verifier 1M-insn complexity budget"
Cohesion: 0.20
Nodes (8): PR vs base verifier cost comparison, Verifier gate workflow, bpf_loop() helper, tools/bpfstats verifier cost tool, static __noinline BPF-to-BPF helpers with args struct, rename-over-guarded-file bypass, xres_move cross-resource move POC, Kernel floor (5.8 / 5.10 / 5.17)

### Community 98 - "Inode-based (dev:ino) executable identity"
Cohesion: 0.18
Nodes (7): CheckBPFLSM preflight, internal/infrastructure shared eBPF kernel, Ports & adapters layering (cmd -> usecase -> repository -> engines), jit_provenance ledger (guard_jit_origin), daemon --check / --verifier-only preflight, guard mode (LSM whitelist/blacklist), network-guard mode

### Community 99 - "CandidateDir"
Cohesion: 0.08
Nodes (23): entryWriters(), reserveChain(), symlinkTargetGlob(), underResource(), globBuilder, ParseGlobName(), TestParseGlobName(), expandPlaceholders() (+15 more)

### Community 100 - "policyEditor"
Cohesion: 0.12
Nodes (12): netGuardTab, policyEditor, policyHooks, policyTab, newPolicyTab(), verbFor(), eventsOrAll(), netWarnings() (+4 more)

### Community 101 - "CopyTreeWithProgress"
Cohesion: 0.15
Nodes (9): migrationTarget, writeTempWithMetadata(), classifyMigrationTarget(), readClassifiedFile(), stampCopyMetadata(), CopyTreeWithProgress(), TestCopyTreeProgress(), treeSize() (+1 more)

### Community 102 - "Well-known bypass classes"
Cohesion: 0.24
Nodes (9): btrfs copy-ioctl gate (file_ioctl / file_ioctl_compat), Well-known bypass classes, btrfs_search POC, io_uring POC, open_by_handle_at POC, process_vm_readv POC (remaining monitor gap), raw_block_device POC, statonly stat/statx metadata leak (+1 more)

### Community 103 - "NewDumpHook"
Cohesion: 0.24
Nodes (11): applyVerbosity(), VerboseLevel, BridgeStdLog(), NewDumpHook(), TestBridgeStdLog(), TestDefaultVerboseMirrorsCurrentDisplay(), TestDumpHookFileIsSynced(), TestDumpHookFiltersByVerbosity() (+3 more)

### Community 104 - "jit_provenance.c"
Cohesion: 0.47
Nodes (8): copy_into(), main(), make_temp_copy(), mode_hold(), mode_memfd(), mode_passfd(), mode_self(), report_dlopen()

### Community 105 - "inspector_probe.c"
Cohesion: 0.56
Nodes (8): check_attach(), check_exe(), check_mem(), check_root(), check_vmread(), main(), report(), traced_exec()

### Community 106 - "ptrace_race.c"
Cohesion: 0.31
Nodes (5): main(), peek_buffer(), run_tracer(), run_victim(), wait_for_file()

### Community 107 - "downloadFile"
Cohesion: 0.33
Nodes (4): downloadFile(), TestDownloadFileHTTPError(), TestDownloadFileMode0700(), progressReader

### Community 108 - "Build & sign release assets (reusable workflow)"
Cohesion: 0.31
Nodes (8): Build & sign release assets (reusable workflow), app-listener-gui best-effort build, Promote to stable release workflow, Release workflow (pre-release per push to main), Release promote flow, make check-compatibility, scripts/install.sh one-line installer, update self-update subcommand

### Community 109 - "installPage"
Cohesion: 0.10
Nodes (14): binSource, installAction, installedBinary, InstallMode, installPage, InstallRequest, latestRelease, Options (+6 more)

### Community 110 - "KernelDev"
Cohesion: 0.09
Nodes (30): BinaryStat, isFuseType(), parseMajorMinor(), TestParseMajorMinorMatchesStatEncoding(), TestVouchedDevs(), TestVouchedDevsFailsClosed(), vouchedDevs(), binaryStatOf() (+22 more)

### Community 111 - "mountTable"
Cohesion: 0.24
Nodes (7): MountInfo, MountByID(), parseMountDev(), ParseMountinfo(), PathUnder(), ReadMountinfo(), mountTable

### Community 112 - "NetworkGuardRepository"
Cohesion: 0.12
Nodes (15): NetworkGuardRepository, NewGuardUseCase(), NetworkGuardUseCase, NewNetworkGuardUseCase(), NewNetworkMonitorUseCase(), newFakeMonitorRepo(), newFakeNetworkGuardRepo(), newFakeNetworkMonitorRepo() (+7 more)

### Community 113 - "app-listener Terminal Demo Recording"
Cohesion: 0.39
Nodes (8): BPF Build Pipeline (bpf2go + go build), CLI Subcommand Help (guard, install, monitor, network-guard, network-monitor), app-listener Terminal Demo Recording, Guard TUI Blacklist Mode (eBPF LSM), Guard TUI Whitelist Mode (eBPF LSM), File System Monitor TUI (eBPF), Startup Banner and eBPF Availability Check, TUI Resource Footer (Mem / CPU / Evt/s)

### Community 115 - "shQuote"
Cohesion: 0.16
Nodes (4): IntegrationSuite, IntegrationSuite, shQuote(), IntegrationSuite

### Community 116 - "changelog.go"
Cohesion: 0.15
Nodes (25): changelogText(), cumulativeChangelog(), newReleaseWindow(), releasesSince(), compareStableVersions(), fetchReleases(), filterChannel(), findUpdate() (+17 more)

### Community 119 - "Regenerate eBPF bindings in pinned builder image"
Cohesion: 0.33
Nodes (5): CI workflow (lint + test on non-draft PRs), CI lint job (make lint), CI test job (make test), golangci-lint v2 configuration, Issue #45 (clang 14 guard_path_rename over budget)

### Community 120 - "serveWebSocket"
Cohesion: 0.21
Nodes (8): configureWebSocket(), sameOrigin(), serveWebSocket(), TestSameOrigin(), TestSnapshotHubSingleViewerAndCoalescing(), writeSnapshot(), writeSnapshotLoop(), snapshotHub

### Community 122 - "WebSocket /ws client"
Cohesion: 0.29
Nodes (3): AppListener Shared Session page, WebSocket /ws client, --serve WebSocket mirroring

### Community 124 - "launch_probe.c"
Cohesion: 0.53
Nodes (4): copy(), is_verb(), main(), usage()

### Community 125 - "install.sh"
Cohesion: 0.57
Nodes (6): die(), download(), info(), install.sh script, usage(), warn()

### Community 126 - "mc_tag.h"
Cohesion: 0.27
Nodes (9): mc_attest(), mc_basename(), mc_leader_start(), mc_stamp_fork(), mc_stamp_free(), mc_tag_bits(), netg_exec_applet(), netm_exec_applet() (+1 more)

### Community 127 - ".FilterExistingWhitelist"
Cohesion: 0.08
Nodes (26): CleanupStalePins(), genOfPinFile(), TestCleanupStalePins(), TestCleanupStalePinsMissingBase(), TestGenOfPinFile(), TestRuntimeClassMarkers(), confSafeMatch(), homeMatchConfined() (+18 more)

### Community 128 - "make test-integration (rootful Docker suite)"
Cohesion: 0.40
Nodes (4): make test-integration (rootful Docker suite), app-listener docker-compose service, ptrace_race POC, Docker usage (docker compose)

### Community 130 - "supersede_probe.c"
Cohesion: 0.90
Nodes (4): answer(), main(), read_secret(), run_child()

### Community 131 - "newTestVetter"
Cohesion: 0.19
Nodes (18): newTestVetter(), testBinary(), TestVetterBootstrapRecordsThenChecks(), TestVetterLinkKeyedOnLink(), TestVetterRecordLiveUnderLinkLines(), TestVetterRefusesNewLineAfterBootstrap(), TestVetterResolverOnlyApproved(), TestVetterSeesInPlaceRewriteThroughSharedHashes() (+10 more)

### Community 132 - "trace-app-libs.sh"
Cohesion: 0.60
Nodes (3): conf_directive_values(), list_conf_binaries(), trace-app-libs.sh script

### Community 134 - "configedit_test.go"
Cohesion: 0.47
Nodes (9): newTestConfigEditor(), readFile(), TestConfigEditorApply(), TestConfigEditorRefusals(), TestConfigEditorRefusesNewMulticallLine(), TestConfigEditorRestoresOnFailedReload(), TestRunConfigEditOverTheSocket(), TestRunConfigEditRefusedDuringALiveSession() (+1 more)

### Community 135 - "Bun TMPDIR redirect with reserved .bun-* name"
Cohesion: 0.50
Nodes (4): lib_probe reserved library POC, Bun TMPDIR redirect with reserved .bun-* name, install --diff-catalog, install TUI wizard

### Community 142 - "TempBinary"
Cohesion: 0.12
Nodes (16): TempBinary, TemporaryGrant, TempRule, planTempBlock(), applyTempGrants(), closeTempBinaries(), TemporaryAccess, logTempGrant() (+8 more)

### Community 143 - "NewEventFanout"
Cohesion: 0.24
Nodes (5): NewEventFanout(), TestEventFanoutDropOldestUnderPressure(), TestEventFanoutDuplicatesInOrderAndCloses(), TestEventFanoutStopClosesOutputs(), EventFanout

### Community 144 - "SanitizeText"
Cohesion: 0.15
Nodes (10): multicallAdmissible(), ownerPaths(), runHeadless(), runHeadless(), binaryVetter, vetDecision, launchKey(), TestLaunchRulesFitTheKey() (+2 more)

### Community 145 - "hashFileShared"
Cohesion: 0.21
Nodes (9): binaryVerifyState, HashStore, sharedHashEntry, claimHash(), hashFileShared(), hashPathShared(), publishHash(), settled() (+1 more)

### Community 146 - "catalogRefresher"
Cohesion: 0.16
Nodes (7): drain(), readEvents(), swapFile(), TestDebouncerCoalescesAndBoundsWait(), TestDebouncerMinGap(), catalogRefresher, debouncer

### Community 147 - "net_tester/main.go"
Cohesion: 0.21
Nodes (7): connFD(), tcpClient(), tcpServer(), tcpServerDelayed(), udpRecvLoop(), udpSendLoop(), waitForMarker()

### Community 148 - "edit-protected"
Cohesion: 0.25
Nodes (8): edit-protected, Editing the daemon config, Editor keys, Exit audit, Non-interactive writes, Offline and live mode, Temporary access for other programs: `--forward`, The password and the control socket

### Community 149 - "Monitor"
Cohesion: 0.15
Nodes (3): NewMonitor(), Monitor, watchedInode

### Community 150 - "hashCache"
Cohesion: 0.18
Nodes (7): newHashCache(), ds(), TestHashCacheFlushKeepsOnlyThisRunsRows(), hashCache, fakeHashStore, DurableStat, HashCache

### Community 152 - "daemon: permanent protection with encryption at rest"
Cohesion: 0.17
Nodes (12): Boot ordering, daemon: permanent protection with encryption at rest, Filesystem-wide gates, Flags, Lifecycle, Process inspectors and taint, Quiet metadata denials, Related (+4 more)

### Community 153 - "ConfinedEntry"
Cohesion: 0.24
Nodes (21): ConfinedEntry(), HashStats(), hashedReused(), resetSharedHashes(), TestConcurrentCallersHashOnce(), TestConfinedEntryHashesOncePerFingerprint(), TestNewVerifyStateSeededOnlyByMatchingHash(), TestPersistentStoreFilledAfterSyncAndHit() (+13 more)

### Community 154 - "EventType"
Cohesion: 0.08
Nodes (15): systemRule(), snapEngines, snapMonitor, BpfEvent, EventType, eventUnitTest, FileEvent, parseEvents() (+7 more)

### Community 155 - ".Layout"
Cohesion: 0.44
Nodes (5): adaptiveGrid, gapVBox, cellMin(), rowHeights(), visible()

### Community 156 - "RunWith"
Cohesion: 0.11
Nodes (19): needsBinary(), RunWith(), TestNeedsBinary(), TestRunWithRefusesDeployWithoutBinary(), TestRunWithRefusesNonHeadlessBinary(), applyLiveRefresh(), runUpdateCatalogOnly(), softenAutomatedRefreshErr() (+11 more)

### Community 157 - "newTestTarget"
Cohesion: 0.11
Nodes (21): FileEngines, NewFiles(), newTestTarget(), since(), TestAddRefusesOverlap(), TestApplyFromBlacklistFailureKeepsOldGuard(), TestApplyFromBlacklistOverlaps(), TestApplyFromWhitelistFailureReportsUnguarded() (+13 more)

### Community 158 - "Installation"
Cohesion: 0.18
Nodes (11): Build from source, Bun-based apps (opencode), Docker, fscrypt prerequisites, Install flags, Installation, One-line installer, SSH agent (+3 more)

### Community 159 - "backend_test.go"
Cohesion: 0.31
Nodes (17): call(), editingBackend(), mustRead(), mustWrite(), newTestBackend(), TestBackendAuthRefusedOffline(), TestBackendAuthTimeoutBounds(), TestBackendChownOnlyLoginUsers() (+9 more)

### Community 161 - "writeWithin"
Cohesion: 0.22
Nodes (9): atomicWriteAt(), descendNoFollow(), refuseExistingSymlink(), TestWriteWithin_PlainFile(), TestWriteWithin_RefusesSymlinkDirComponent(), TestWriteWithin_RefusesSymlinkTarget(), TestWriteWithin_RejectsEscape(), writeAndSync() (+1 more)

### Community 162 - "golang.org/x/sys/unix.Stat_t"
Cohesion: 0.20
Nodes (11): placedWalk, confineBelow(), inUserHome(), GlobSystemPlaced(), placedDir(), placedMatches(), readlinkFd(), rootOwnedStat() (+3 more)

### Community 163 - "main_test.go"
Cohesion: 0.20
Nodes (5): infraContainerPaths(), TestIntegrationSuite(), TestMain(), helperChild(), TestMain()

### Community 164 - "pinstate.go"
Cohesion: 0.27
Nodes (15): pinOwnerLikelyAlive(), readPinState(), recoverPinState(), pinStateInode(), TestPinOwnerLikelyAlive(), TestPinOwnerLikelyAliveDeadPID(), TestPinStateRoundTrip(), TestReadPinStateMalformed() (+7 more)

### Community 165 - "feed.go"
Cohesion: 0.06
Nodes (44): actionBar, cell, column, counterStrip, dataTable, eventFeed, feedAction, feedSpec (+36 more)

### Community 166 - "shell"
Cohesion: 0.06
Nodes (31): pageDef, pageID, resizeWatch, shell, sizedDialog, sysStatus, newLogo(), newShell() (+23 more)

### Community 167 - "integrationtests/networkguard_test.go"
Cohesion: 0.27
Nodes (5): IntegrationSuite, guardNetTypesForComm(), netGuardBlockedEventCount(), netGuardHasBlockedEvent(), netGuardTail()

### Community 168 - "WithBottomBar"
Cohesion: 0.29
Nodes (5): TestProgressReader(), bottomBar, newBottomBar(), renderBar(), WithBottomBar()

### Community 170 - "lockRootRecovering"
Cohesion: 0.18
Nodes (13): lockOneRoot(), lockRootRecovering(), relockStaleVaults(), ensurePinStateFilePlaceholder(), ensurePlaceholder(), bootstrapSelfProtectPlaceholders(), ensureHashFilePlaceholder(), readTempJournal() (+5 more)

### Community 172 - "Limitations"
Cohesion: 0.22
Nodes (9): Already-running processes, Binaries you can write, Event masks are not confinement, Interpreters and JITs, Limitations, Project status, Root, Side effects to expect with the daemon (+1 more)

### Community 173 - "Troubleshooting"
Cohesion: 0.22
Nodes (9): A program I use is denied, After a package update replaced a binary, An app breaks and nothing is logged, Backups, Check the guard by hand, Encryption key errors, Is the daemon running and enforcing?, Snapshots or disk backups stopped working (+1 more)

### Community 176 - "Daemon configuration"
Cohesion: 0.25
Nodes (8): Daemon configuration, Daemon settings: `[settings]`, Generic Electron apps: `[electron_apps]`, Library trust: `[libraries "<name>"]`, Per-binary event masks, Reloading, Several guarded trees, one vault, Watch sections

### Community 177 - "Active"
Cohesion: 0.16
Nodes (12): showChangelog(), confirmUpdate(), isTerminal(), RunWith(), TestConfirmUpdateRoutesToActiveUI(), TestConfirmUpdateWithoutTerminalOrUIRefuses(), TestRunWithRefusesConcurrentRun(), cancelUI (+4 more)

### Community 178 - "batchInodeRows"
Cohesion: 0.19
Nodes (12): markReady(), notifySystemdReady(), inodeBatchFn, appendInodeRow(), batchInodeRows(), inodeRowsOf(), fakeBatches(), TestBatchInodeRows() (+4 more)

### Community 179 - "CheckBPFLSM"
Cohesion: 0.17
Nodes (13): CheckBPFLSM(), confirmMasterKeyOverwrite(), runBPFCheck(), runGenKey(), runOneShotMode(), ensureInstalledBinary(), ensureMasterKey(), prepareInstallation() (+5 more)

### Community 180 - "Compatibility"
Cohesion: 0.25
Nodes (8): Architectures, BPF-LSM must be active, not just compiled in, Check your host first, Compatibility, Distributions, Kernel floors per mode, Low-end hardware, Mandatory and best-effort hooks

### Community 181 - "How it works"
Cohesion: 0.25
Nodes (8): eBPF at three levels, Fail closed, Global flags, How it works, Identity is the executable's inode, Output, Watching from a browser: `--serve`, Where it sits: DAC, MAC and app-listener

### Community 182 - "ResolveTempBinary"
Cohesion: 0.18
Nodes (8): confirmUserPlaced(), forwardCheck(), binaryStatOf(), loadPinnedExeMaps(), ResolveTempBinary(), stripPinnedAllow(), StripPinnedTempAllows(), TestResolveTempBinaryClassesRuntime()

### Community 183 - "editControlSession"
Cohesion: 0.30
Nodes (5): readClientLines(), streamEvents(), TestEditControlSessionEndIdempotent(), editControlSession, grantRequest

### Community 185 - "Parse"
Cohesion: 0.33
Nodes (8): dialLiveSession(), applyEdit(), checkAndApply(), editAgain(), fetchConfig(), putConfig(), runEditConfig(), Parse()

### Community 186 - "Binary trust"
Cohesion: 0.29
Nodes (7): Binary ledger, Binary replacement, Binary trust, Hash cache, Launch scan, Multicall binaries, trust-binaries

### Community 187 - "guard: block file access by program"
Cohesion: 0.33
Nodes (6): Executables inside the guarded tree, Flags, guard: block file access by program, Multicall binaries, Replaced binaries, Whitelist and blacklist

### Community 188 - "network-guard: block network access by program"
Cohesion: 0.33
Nodes (6): Code integrity (`-w`), Flags, Multicall binaries, network-guard: block network access by program, Replaced binaries, What code integrity cannot stop

### Community 189 - "lockAndDeprovision"
Cohesion: 0.21
Nodes (8): deprovisionKind, classifyDeprovision(), classifyDeprovisionErr(), applyRawKeyPolicy(), isDeprovisionBusy(), isDeprovisionMissing(), lockAndDeprovision(), modifiedContextWithSource()

### Community 190 - "Development"
Cohesion: 0.33
Nodes (6): Demo GIFs, Development, Docker, Makefile targets, Profiling, Tests

### Community 191 - "app-listener documentation"
Cohesion: 0.40
Nodes (5): app-listener documentation, Getting started, Modes, Operating, The daemon in depth

### Community 192 - "Self-updating apps and the catalog"
Cohesion: 0.50
Nodes (4): Catalog refresh, Self-updating apps, Self-updating apps and the catalog, The catalog

### Community 193 - "network-monitor: see what a program talks to"
Cohesion: 0.50
Nodes (4): Events, Flags, Multicall binaries, network-monitor: see what a program talks to

### Community 195 - "github.com/stretchr/testify/suite.Suite"
Cohesion: 0.18
Nodes (4): harnessSuite, mountinfoUnitTest, unescapeMountinfo(), editorHarness

### Community 196 - "bufio.Reader"
Cohesion: 0.22
Nodes (10): TestReadGrantRequestConfig(), readAuthRequest(), readForwardRequest(), readGrantRequest(), readPathLines(), TestReadAuthRequest(), TestReadGrantRequestForward(), TestReadGrantRequestSelect() (+2 more)

### Community 197 - "monitor: see what touches a path"
Cohesion: 0.67
Nodes (3): Events, Flags, monitor: see what touches a path

### Community 199 - "pathCache"
Cohesion: 0.39
Nodes (4): evalSymlinksOrEmpty(), newPathCache(), pathCache, pathCacheEntry

### Community 200 - ".eventTypeAllowed"
Cohesion: 0.29
Nodes (3): TestEventTypeAllowedConcurrent(), TestSetEventTypesBeforeStartNoRace(), TestSetEventTypesRace()

### Community 201 - "Execute"
Cohesion: 0.18
Nodes (10): installHooks(), releaseBuild(), Requested(), Run(), TestRequested(), TestSubcommandOutranksGUI(), Execute(), initLogger() (+2 more)

### Community 202 - "decryptStep"
Cohesion: 0.20
Nodes (9): askDecrypt(), cleanOrphanedMetadata(), decryptStep(), pickDirsToDecrypt(), CleanOrphanedMetadata(), LivePolicyDescriptors(), TestLivePolicyDescriptorsTmpfs(), VerifyEncryptedKeys() (+1 more)

### Community 204 - "WithSuspectsFrom"
Cohesion: 0.36
Nodes (4): WithSuspectsFrom(), RealNetEngines(), netguardHarness, Option

### Community 205 - "copyRegularAt"
Cohesion: 0.31
Nodes (10): copyDirContents(), copyEntry(), copyRegularAt(), copySymlinkAt(), copyTreeRoot(), fchownFromInfo(), futimesFromInfo(), readDirNames() (+2 more)

### Community 206 - "Desktop GUI"
Cohesion: 0.17
Nodes (12): Actions, Desktop GUI, File and folder pickers, Files, Home, Install, Network, Protected (+4 more)

### Community 208 - "ParseEventsFlag"
Cohesion: 0.29
Nodes (6): ParseEventsFlag(), TestValidateFlagsEditConfig(), validateFlags(), TestValidateForwardFlags(), validateForwardFlags(), validateForwardRule()

### Community 209 - "runLockdown"
Cohesion: 0.33
Nodes (6): configureDaemonLogging(), loadDaemonConfig(), resolveConfigPath(), runLockdown(), TestLoadDaemonConfigEmptyResourcesIsCriticalStartup(), TestLoadDaemonConfigMissingFileIsCriticalStartup()

### Community 210 - "parseResize"
Cohesion: 0.33
Nodes (4): decodeEndsAtEOF(), parseResize(), readResizeLoop(), TestParseResize()

### Community 211 - "checkRunningBinaryMatchesInstalled"
Cohesion: 0.40
Nodes (5): checkInstalledMatchesSource(), checkRunningBinaryMatchesInstalled(), fileSHA256(), sameContent(), TestCheckRunningBinaryMatchesInstalled()

### Community 212 - "TestSuperblockDevMatchesStatOffBtrfs"
Cohesion: 0.50
Nodes (4): FirstOnBtrfs(), OnBtrfs(), SuperblockDevPath(), TestSuperblockDevMatchesStatOffBtrfs()

## Knowledge Gaps
- **32 isolated node(s):** `github.com/Virgula0/app-listener`, `IntegrationSuite`, `open_by_handle_at POC`, `statonly stat/statx metadata leak`, `app-listener-gui best-effort build` (+27 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 545 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **35 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `edit-protected` connect `edit-protected` to `README.md`?**
  _High betweenness centrality (0.036) - this node is a cross-community bridge._
- **What connects `github.com/Virgula0/app-listener`, `IntegrationSuite`, `open_by_handle_at POC` to the rest of the system?**
  _32 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `daemon/daemon.go` be split into smaller, more focused modules?**
  _Cohesion score 0.09209669867878635 - nodes in this community are weakly interconnected._
- **Why does `Editor keys` connect `edit-protected` to `liveSession`, `runUninstall`?**
  _High betweenness centrality (0.035) - this node is a cross-community bridge._
- **Should `testing.T` be split into smaller, more focused modules?**
  _Cohesion score 0.019955849889624726 - nodes in this community are weakly interconnected._
- **Why does `absPath()` connect `absPath` to `IntegrationSuite`, `main_test.go`, `integrationtests/networkguard_test.go`, `IntegrationSuite`, `IntegrationSuite`, `.launch`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `shQuote`, `IntegrationSuite`, `IntegrationSuite`, `daemon_multicall_test.go`, `IntegrationSuite`?**
  _High betweenness centrality (0.031) - this node is a cross-community bridge._
- **Should `IntegrationSuite` be split into smaller, more focused modules?**
  _Cohesion score 0.08957708957708958 - nodes in this community are weakly interconnected._